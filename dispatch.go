package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	shield "github.com/crasm/shield/lib"
)

const FileExtension = `.shd`

// Figures out input and output files and calls the appropiate shield library
// functions on them.
func dispatch(in, out string) error {
	var err error

	inFile := os.Stdin
	outFile := os.Stdout

	explicitIn := in != Stdio && in != ""
	explicitOut := out != Stdio && out != ""

	if explicitIn {
		inFile, err = os.Open(in)
		defer inFile.Close()
		if err != nil {
			return err
		}
	}

	implicitOut := explicitIn && out == ""

	if explicitOut {
		outFile, err = createFile(out, opt.Force)
		if err != nil {
			// TODO: If we can't truncate, then try to append. This
			// should make using /dev/stderr and /dev/stdout work as
			// expected.
			return err
		}
	}

	switch {
	case opt.Create:
		if implicitOut {
			inferred := fmt.Sprint(in, FileExtension)
			outFile, err = createFile(inferred, opt.Force)
			defer outFile.Close()
			if err != nil {
				return err
			}
		}

		bytes, err := bitsToBytes(opt.Size)
		if err != nil {
			return err
		}

		if outFile == os.Stdout {
			_, err = shield.WrapBuffered(inFile, outFile, bytes)
		} else {
			_, err = shield.Wrap(inFile, outFile, bytes)
		}

	case opt.Extract:
		if implicitOut {
			inferred := strings.TrimSuffix(in, FileExtension)
			outFile, err = createFile(inferred, opt.Force)
			defer outFile.Close()
			if err != nil {
				die(err)
			}
		}
		_, err = shield.Unwrap(inFile, outFile)

	case opt.Verify:
		var shd *shield.UnwrappedShield
		shd, err = shield.Unwrap(inFile, ioutil.Discard)
		fmt.Fprintf(outFile, "claim:  %v\nactual: %v\n",
			hex.EncodeToString(shd.Claim),
			hex.EncodeToString(shd.Actual))

	case opt.Dump:
		err = shield.DumpHeader(inFile, outFile)
	default:
		panic("No command specified")
	}

	return err
}
