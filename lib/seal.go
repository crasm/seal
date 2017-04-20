// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package seal

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const Magic = `SL%v`
const Version = 0

const IdentLen = len(`SL%v0`)

const DefaultSealBits = 512

var ErrSealBroken = errors.New("seal: claim did not validate against content")
var ErrBadSignatureLength = errors.New("seal: signature length is invalid")

const maxBytes = sha512.Size

// Seal is the information stored in the seal header.
type Seal struct {
	Magic            string
	Version          int
	ClaimedSignature []byte
}

// UnwrappedSeal extends Seal to provide the calculated signature of the
// sealed content.
type UnwrappedSeal struct {
	Seal
	CalculatedSignature []byte
}

func (sl *Seal) Bytes() []byte {
	return []byte(sl.String())
}

func (sl *Seal) String() string {
	return fmt.Sprintf("%s%d{%s}\n", sl.Magic, sl.Version,
		hex.EncodeToString(sl.ClaimedSignature))
}

// Generate a Seal for the io.Reader with the default number of bits.
func Generate(in io.Reader) (*Seal, error) {
	return GenerateBits(in, DefaultSealBits)
}

func GenerateBits(in io.Reader, bits int) (*Seal, error) {
	sigLen := bitsToBytes(bits)
	if sigLen == -1 {
		return nil, ErrBadSignatureLength
	}

	calc, err := sum(in)

	return &Seal{
		Magic:            Magic,
		Version:          Version,
		ClaimedSignature: calc[:sigLen],
	}, err
}

// Wrap the contents of `in` with a Seal header, and write the full Seal
// file to `out`. Uses the default number of bits.
func Wrap(in io.Reader, out io.WriteSeeker) (*Seal, error) {
	return WrapBits(in, out, DefaultSealBits)
}

func WrapBits(in io.Reader, out io.WriteSeeker, bits int) (*Seal, error) {
	var err error

	sigLen := bitsToBytes(bits)
	if sigLen == -1 {
		return nil, ErrBadSignatureLength
	}

	contentOffset := headerLen(sigLen)

	_, err = out.Seek(int64(contentOffset), 0)
	if err != nil {
		return nil, err
	}

	calc, err := teesum(in, out)
	if err != nil {
		return nil, err
	}

	sl := &Seal{
		Magic:            Magic,
		Version:          Version,
		ClaimedSignature: calc[:sigLen],
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	_, err = out.Write(sl.Bytes())

	return sl, err
}

// Same as Wrap, but uses a temporary file to buffer the output because
// `out` is not seekable.
func WrapBuffered(in io.Reader, out io.Writer) (*Seal, error) {
	return WrapBufferedBits(in, out, DefaultSealBits)
}

func WrapBufferedBits(in io.Reader, out io.Writer, bits int) (*Seal, error) {
	tmp, err := ioutil.TempFile("", "seal")
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err != nil {
		return nil, err
	}

	// Do the actual wrapping, but output to a temporary file.
	sl, err := WrapBits(in, tmp, bits)
	if err != nil {
		return sl, err
	}

	_, err = tmp.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	outwr := bufio.NewWriter(out)
	_, err = outwr.ReadFrom(tmp)
	if err != nil {
		return sl, err
	}

	return sl, outwr.Flush()
}

func Unwrap(in io.Reader, out io.Writer) (*UnwrappedSeal, error) {

	bufIn := bufio.NewReader(in)

	s, err := parseHeader(bufIn)
	if err != nil {
		return nil, err
	}

	sl := &UnwrappedSeal{Seal: *s}

	calcSig, err := teesum(bufIn, out)
	if err != nil {
		return sl, err
	}

	sl.CalculatedSignature = calcSig[:len(sl.ClaimedSignature)]

	if !bytes.Equal(sl.ClaimedSignature, sl.CalculatedSignature) {
		err = ErrSealBroken
	}

	return sl, err
}

// Dump the raw seal header.
func DumpHeader(in io.Reader, out io.Writer) error {
	// TODO: Make this safer. Search for `}\n` or `}\r\n`
	line, err := bufio.NewReader(in).ReadSlice('\n')
	if err != nil {
		return err
	}

	_, err = out.Write(line)
	return err
}

func bitsToBytes(bits int) int {
	bytes := bits / 8
	if bits <= 0 || bytes*8 != bits || bytes > maxBytes {
		return -1
	}
	return bytes
}

func sum(in io.Reader) ([]byte, error) {
	digester := sha512.New()
	_, err := bufio.NewReader(in).WriteTo(digester)
	return digester.Sum(nil), err
}

// Take the hash of the data from in, write it to out, and return the hash.
func teesum(in io.Reader, out io.Writer) ([]byte, error) {
	digester := sha512.New()

	tee := bufio.NewReader(io.TeeReader(in, out))
	_, err := tee.WriteTo(digester)
	if err != nil {
		return nil, err
	}

	return digester.Sum(nil), nil
}
