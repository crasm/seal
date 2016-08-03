package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	shield "github.com/crasm/shield/lib"
)

// Figures out input and output files and calls the appropiate shield library
// functions on them.
func dispatch(cmd Command, in, out *os.File) error {

	// TODO: Move this validation code somewhere else. (Library?)
	bytes, err := bitsToBytes(opt.Size)
	if err != nil {
		return err
	}

	digester := shield.NewDigesterSha512(bytes)

	switch cmd {
	case Create:
		if out.Name() == os.Stdout.Name() {
			_, err = digester.WrapBuffered(in, out)
		} else {
			_, err = digester.Wrap(in, out)
		}

	case Extract:
		_, err = digester.Unwrap(in, out)

	case Verify:
		var shd *shield.UnwrappedShield
		shd, err = digester.Unwrap(in, ioutil.Discard)
		fmt.Fprintf(out, "claim:  %v\nactual: %v\n",
			hex.EncodeToString(shd.Claim),
			hex.EncodeToString(shd.Actual))

	case Dump:
		err = shield.DumpHeader(in, out)
	default:
		panic("no command specified")
	}

	return err
}
