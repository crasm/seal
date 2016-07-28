package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"fmt"
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

		if outFile == os.Stdout {
			err = shield.WrapBuffered(inFile, outFile)
		} else {
			err = shield.Wrap(inFile, outFile)
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
		err, _, _ = shield.Unwrap(inFile, outFile)
	case opt.Dump:
		fallthrough
	default:
		die("Dump command is not supported yet.")
	}

	return err
}
