// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package seal

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

// Returns the number of bytes in a header of size "size".
func headerSize(bytes int) (int, error) {
	if bytes <= 0 || bytes > maxBytes {
		// TODO provide bounds
		return -1, fmt.Errorf("seal: Invalid header length: %v", bytes*8)
	}

	// Hex encoding requires two bytes per byte, so we multiply by 2.
	return IdentLen + len("{}\n") + bytes*2, nil
}

// Parses the header of a seal file. Does not read beyond the
// header.
func parseHeader(in io.Reader) (*Seal, error) {
	var err error

	limit, _ := headerSize(maxBytes)
	header := make([]byte, limit)

	// We read byte-by-byte because we don't know how long the header is.
	var i int
	for i = 0; i < limit; i++ {
		_, err := in.Read(header[i : i+1])
		if err != nil && err != io.EOF {
			fmt.Errorf("seal: Error reading header: %v", err)
		}

		// TODO: Should be terminated by `}\n` or `}\r\n`
		if header[i] == '\n' {
			break
		}
	}
	header = header[:i+1] // Reslice so we don't have trailing zeros.

	sl := &Seal{}

	sl.Magic = string(header[:len(Magic)]) // example: `SL%v`
	if sl.Magic != Magic {
		return sl, fmt.Errorf("Got wrong magic number: %v", sl.Magic)
	}

	uint64Version, err := strconv.ParseUint(
		string(header[len(Magic)]), 10, 8) // example: `0`
	if err != nil {
		return sl, fmt.Errorf("Couldn't parse version: %v", err)
	}

	sl.Version = int(uint64Version)
	if sl.Version != Version {
		return nil, fmt.Errorf("Unsupported version: %v", sl.Version)
	}

	hexClaim := strings.TrimFunc(string(header[IdentLen:]),
		func(r rune) bool { return unicode.IsSpace(r) || r == '{' || r == '}' },
	)
	sl.Claim, err = hex.DecodeString(hexClaim)

	return sl, err
}

func createHeader(claim []byte) string {
	hexclaim := hex.EncodeToString(claim)
	return fmt.Sprintf("%s%d{%s}\n", Magic, Version, hexclaim)
}
