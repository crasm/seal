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

// Determine the proper input and output files. If outFile can't be
// determined (is implicit), then nil is returned for outFile.
func determineInputOutput(inArg, outArg string) (inFile, outFile *os.File, err error) {
	inFile = os.Stdin
	outFile = os.Stdout

	explicitIn := inArg != Stdio && inArg != ""
	explicitOut := outArg != Stdio && outArg != ""

	if explicitIn {
		inFile, err = os.Open(inArg)
		if err != nil {
			return
		}
	}

	if explicitOut {
		outFile, err = createFile(outArg, opt.Force)
		if err != nil {
			// TODO: If we can't truncate, then try to append. This
			// should make using /dev/stderr and /dev/stdout work as
			// expected.
			return
		}
	}

	implicitOut := explicitIn && outArg == ""
	if implicitOut {
		outFile = nil
	}

	return inFile, outFile, nil
}

// Figures out input and output files and calls the appropiate shield library
// functions on them.
func dispatch(in, out *os.File) error {

	implicitOut := out == nil

	// TODO: Move this validation code somewhere else. (Library?)
	bytes, err := bitsToBytes(opt.Size)
	if err != nil {
		return err
	}

	digester := shield.NewDigesterSha512(bytes)

	switch {
	case opt.Create:
		if implicitOut {
			inferred := fmt.Sprint(in.Name(), FileExtension)
			logger.Debug("Inferring output file as %v", inferred)
			out, err = createFile(inferred, opt.Force)
			defer out.Close()
			if err != nil {
				return err
			}
		}

		if out == os.Stdout {
			_, err = digester.WrapBuffered(in, out)
		} else {
			_, err = digester.Wrap(in, out)
		}

	case opt.Extract:
		if implicitOut {
			inferred := strings.TrimSuffix(in.Name(), FileExtension)
			logger.Debug("Inferring output file as %s", inferred)
			out, err = createFile(inferred, opt.Force)
			defer out.Close()
			if err != nil {
				die(err)
			}
		}
		_, err = digester.Unwrap(in, out)

	case opt.Verify:
		var shd *shield.UnwrappedShield
		shd, err = digester.Unwrap(in, ioutil.Discard)
		fmt.Fprintf(out, "claim:  %v\nactual: %v\n",
			hex.EncodeToString(shd.Claim),
			hex.EncodeToString(shd.Actual))

	case opt.Dump:
		err = shield.DumpHeader(in, out)
	default:
		panic("No command specified")
	}

	return err
}
