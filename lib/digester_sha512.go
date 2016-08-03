// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package shield

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
)

type digesterSha512 struct {
	hash  hash.Hash
	trunc int
}

// TODO NEXT: This solves the problem of "where do I put trunclen?" But
// now, need to properly implement the rest of it.
func NewDigesterSha512(trunc int) Digester {
	return &digesterSha512{hash: sha512.New(), trunc: trunc}
}

// Creates a shield on the data from in. Writes the shield header and the file
// contents to out. In most cases, out should be an *os.File. However, anything
// that supports seeking to the start is supported.
//
// Returns a Shield describing the header of the shield file just created.
// trunc is bytes, not bits
func (d *digesterSha512) Wrap(in io.Reader, out io.WriteSeeker) (*Shield, error) {
	size, err := headerSize(d.trunc)
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

	calc = calc[:d.trunc]

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
func (d *digesterSha512) WrapBuffered(in io.Reader, out io.Writer) (*Shield, error) {
	tmp, err := ioutil.TempFile("", "shield")
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err != nil {
		return nil, err
	}

	// Do the actual wrapping, but output to a temporary file.
	shd, err := d.Wrap(in, tmp)
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
func (d *digesterSha512) Unwrap(in io.Reader, out io.Writer) (*UnwrappedShield, error) {
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
