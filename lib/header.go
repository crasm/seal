// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package shield

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const maxBytes = sha512.Size

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

// Returns the number of bytes in a header of size "size".
func headerSize(bytes int) (int, error) {
	if bytes <= 0 || bytes > maxBytes {
		// TODO provide bounds
		return -1, fmt.Errorf("shield: Invalid header length: %v", bytes*8)
	}

	// Hex encoding requires two bytes per byte, so we multiply by 2.
	return IdentLen + len("{}\n") + bytes*2, nil
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

func createHeader(claim []byte) string {
	hexclaim := hex.EncodeToString(claim)
	return fmt.Sprintf("%sv%d{%s}\n", Magic, Version, hexclaim)
}
