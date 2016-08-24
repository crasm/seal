// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package seal

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"hash"
	"io"
	"io/ioutil"
	"os"
)

type digesterSha512 struct {
	hash  hash.Hash
	trunc int
}

// Creates a new digester with a sha512 hash truncated to trunc bytes.
func NewDigesterSha512(trunc int) Digester {
	return &digesterSha512{hash: sha512.New(), trunc: trunc}
}

// Creates a seal on the data from in. Writes the seal header and the file
// contents to out. In most cases, out should be an *os.File. However, anything
// that supports seeking to the start is supported.
//
// Returns a Seal describing the header of the seal file just created.
// trunc is bytes, not bits
func (d *digesterSha512) Wrap(in io.Reader, out io.WriteSeeker) (*Seal, error) {
	size, err := headerSize(d.trunc)
	if err != nil {
		return nil, err
	}

	_, err = out.Seek(int64(size), 0)
	if err != nil {
		return nil, err
	}

	sl := &Seal{}

	calc, err := teesum(in, out)
	if err != nil {
		return nil, err
	}

	calc = calc[:d.trunc]

	_, err = out.Seek(0, 0)
	if err != nil {
		return sl, err
	}

	header := createHeader(calc)
	_, err = out.Write([]byte(header))
	return sl, err
}

// Creates a seal on the data from in, buffering the input to a temporary
// file.
func (d *digesterSha512) WrapBuffered(in io.Reader, out io.Writer) (*Seal, error) {
	tmp, err := ioutil.TempFile("", "seal")
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err != nil {
		return nil, err
	}

	// Do the actual wrapping, but output to a temporary file.
	sl, err := d.Wrap(in, tmp)
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

// Writes the file contents after the header to out. If the claim does not
// validate (match the actual hash), a non-nil error is returned.
func (d *digesterSha512) Unwrap(in io.Reader, out io.Writer) (*UnwrappedSeal, error) {
	s, err := parseHeader(in)
	if err != nil {
		return nil, err
	}

	sl := &UnwrappedSeal{Seal: *s}

	actual, err := teesum(in, out)
	if err != nil {
		return sl, err
	}

	sl.Actual = actual[:len(sl.Claim)]

	if !bytes.Equal(sl.Claim, sl.Actual) {
		err = ErrSealBroken
	}

	return sl, err
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
