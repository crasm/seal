package shield

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const Magic = `SHD%`
const Version = 0 // Current major version number for the library.

const IdentLen = len(`SHD%v0`)

const maxBytes = sha512.Size

// Shield is the information read from a shield header.
type Shield struct {
	Magic   string
	Version int
	Claim   []byte
}

// Shield contains shield file data.
// TODO: Clarify.
type UnwrappedShield struct {
	Shield
	Actual []byte
}

// Returns the number of bytes in a header of size "size".
func headerSize(bytes int) (int, error) {
	if bytes <= 0 || bytes > maxBytes {
		// TODO provide bounds
		return -1, fmt.Errorf("shield: Invalid header length: %v", bytes*8)
	}

	// Hex encoding requires two bytes per byte, so we multiply by 2.
	return IdentLen + len("{}\n") + bytes*2, nil
}

// Creates a shield on the data from in. Writes the shield header and the file
// contents to out. In most cases, out should be an *os.File. However, anything
// that supports seeking to the start is supported.
//
// Returns a Shield describing the header of the shield file just created.
// trunc is bytes, not bits
func Wrap(in io.Reader, out io.WriteSeeker, trunc int) (*Shield, error) {
	size, err := headerSize(trunc)
	if err != nil {
		return nil, err
	}

	_, err = out.Seek(int64(size), 0)
	if err != nil {
		return nil, err
	}

	shd := &Shield{}

	calc, err := teesum(in, out)
	if err != nil {
		return nil, err
	}

	calc = calc[:trunc]

	_, err = out.Seek(0, 0)
	if err != nil {
		return shd, err
	}

	header := createHeader(calc)
	_, err = out.Write([]byte(header))
	return shd, err
}

// Creates a shield on the data from in, buffering the input to a temporary
// file.
func WrapBuffered(in io.Reader, out io.Writer, trunc int) (*Shield, error) {
	tmp, err := ioutil.TempFile("", "shield")
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err != nil {
		return nil, err
	}

	// Do the actual wrapping, but output to a temporary file.
	shd, err := Wrap(in, tmp, trunc)
	if err != nil {
		return shd, err
	}

	outwr := bufio.NewWriter(out)
	_, err = outwr.ReadFrom(tmp)
	if err != nil {
		return shd, err
	}

	return shd, outwr.Flush()
}

// Writes the file contents after the header to out. If the claim does not
// validate (match the actual hash), a non-nil error is returned.
func Unwrap(in io.Reader, out io.Writer) (*UnwrappedShield, error) {
	s, err := parseHeader(in)
	if err != nil {
		return nil, err
	}

	shd := &UnwrappedShield{Shield: *s}

	actual, err := teesum(in, out)
	if err != nil {
		return shd, err
	}

	shd.Actual = actual[:len(shd.Claim)]

	if !bytes.Equal(shd.Claim, shd.Actual) {
		// TODO: Make this a const error.
		err = fmt.Errorf("shield: claim did not match actual hash")
	}

	return shd, err
}

// Dump the raw shield header.
func DumpHeader(in io.Reader, out io.Writer) error {
	// TODO: Make this safer. Search for `}\n` or `}\r\n`
	line, err := bufio.NewReader(in).ReadSlice('\n')
	if err != nil {
		return err
	}

	_, err = out.Write(line)
	return err
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

func createHeader(claim []byte) string {
	hexclaim := hex.EncodeToString(claim)
	return fmt.Sprintf("%sv%d{%s}\n", Magic, Version, hexclaim)
}

// Parses the header of a shield file. Does not read beyond the
// header.
func parseHeader(in io.Reader) (*Shield, error) {
	var err error

	limit, _ := headerSize(maxBytes)
	header := make([]byte, limit)

	// We read byte-by-byte because we don't know how long the header is.
	var i int
	for i = 0; i < limit; i++ {
		_, err := in.Read(header[i : i+1])
		if err != nil && err != io.EOF {
			fmt.Errorf("shield: Error reading header: %v", err)
		}

		// TODO: Should be terminated by `}\n` or `}\r\n`
		if header[i] == '\n' {
			break
		}
	}
	header = header[:i+1] // Reslice so we don't have trailing zeros.

	shd := &Shield{}

	shd.Magic = string(header[:len(Magic)]) // example: `SHD%`
	if shd.Magic != Magic {
		return shd, fmt.Errorf("Got wrong magic number: %v", shd.Magic)
	}

	byteVersion := header[len(Magic) : len(Magic)+2] // example: `v0`
	uint64Version, err := strconv.ParseUint(
		string(byteVersion[1:]), 10, 8) // example: `0`
	if err != nil {
		return shd, fmt.Errorf("Couldn't parse version: %v", err)
	}

	shd.Version = int(uint64Version)
	if shd.Version != Version {
		return nil, fmt.Errorf("Unsupported version: %v", shd.Version)
	}

	hexClaim := strings.TrimFunc(string(header[IdentLen:]),
		func(r rune) bool { return unicode.IsSpace(r) || r == '{' || r == '}' },
	)
	shd.Claim, err = hex.DecodeString(hexClaim)

	return shd, err
}
