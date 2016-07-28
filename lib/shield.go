package shield

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const Magic = `SHD%`
const Version = 0 // Current major version number for the library.

const IdentLen = len(`SHD%v0`)
const HeaderLen = len(`SHD%v0{e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855}`) + len("\n")

// Creates a shield on the data from in. Writes the shield header and the file
// contents to out. In most cases, out should be an *os.File. However, anything
// that supports seeking to the start is supported.
func Wrap(in io.Reader, out io.WriteSeeker) error {
	_, err := out.Seek(int64(HeaderLen), 0)
	if err != nil {
		return err
	}

	calc, err := teesum(in, out)
	if err != nil {
		return err
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		return err
	}

	header := createHeader(calc)
	_, err = out.Write([]byte(header))
	return err
}

// Creates a shield on the data from in, buffering the input to a temporary
// file.
func WrapBuffered(in io.Reader, out io.Writer) error {
	tmp, err := ioutil.TempFile("", "shield")
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	if err != nil {
		return err
	}

	calc, err := teesum(in, tmp)
	if err != nil {
		return err
	}

	header := createHeader(calc)

	_, err = tmp.Seek(0, 0)
	if err != nil {
		return err
	}

	outwr := bufio.NewWriter(out)
	_, err = outwr.WriteString(header)
	if err != nil {
		return err
	}

	_, err = outwr.ReadFrom(tmp)
	if err != nil {
		return err
	}

	return outwr.Flush()
}

// Writes the file contents after the header to out. If the claim does not
// validate (match the actual hash), a non-nil error is returned.
func Unwrap(in io.Reader, out io.Writer) (err error, claim, actual []byte) {

	claim, err = parseHeader(in)
	if err != nil {
		return err, nil, nil
	}

	actual, err = teesum(in, out)
	if err != nil {
		return err, nil, nil
	}

	if !bytes.Equal(claim, actual) {
		return fmt.Errorf("shield: claim did not match actual hash"), claim, actual
	}

	return nil, claim, actual
}

// Take the hash of the data from in, write it to out, and return the hash.
func teesum(in io.Reader, out io.Writer) ([]byte, error) {
	digester := sha256.New()

	tee := bufio.NewReader(io.TeeReader(in, out))
	_, err := tee.WriteTo(digester)
	if err != nil {
		return nil, err
	}

	return digester.Sum(nil), nil
}

// Pipes a shielded file (sans header) "in" to os.Stdout and verifies the
// contents.
func Pipe(in io.Reader) error {
	return nil
}

// Copies a file from in to out, verifying the shield contents.
func Copy(in, out string) error {
	return nil
}

func createHeader(claim []byte) string {
	hexclaim := hex.EncodeToString(claim)
	return fmt.Sprintf("%sv%d{%s}\n", Magic, Version, hexclaim)
}

func parseHeader(in io.Reader) ([]byte, error) {
	header := make([]byte, HeaderLen)
	n, err := in.Read(header)
	if err != nil && err != io.EOF {
		return nil, err
	} else if n < HeaderLen {
		return nil, fmt.Errorf("shield: Couldn't read a full header.")
	}

	ident := header[:IdentLen]

	if m := string(ident[:len(Magic)]); m != Magic {
		return nil, fmt.Errorf("Got wrong magic number: %s", m)
	}
	if v := string(ident[len(Magic):]); v != fmt.Sprintf("v%d", Version) {
		return nil, fmt.Errorf("Got wrong version: %s", v)
	}

	hexClaim := strings.TrimFunc(string(header[IdentLen:]),
		func(r rune) bool {
			return unicode.IsSpace(r) || r == '{' || r == '}'
		},
	)

	return hex.DecodeString(hexClaim)
}
