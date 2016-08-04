// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	seal "github.com/crasm/seal/lib"
)

// Figures out input and output files and calls the appropiate seal library
// functions on them.
func dispatch(cmd Command, in, out *os.File) error {

	// TODO: Move this validation code somewhere else. (Library?)
	bytes, err := bitsToBytes(opt.Size)
	if err != nil {
		return err
	}

	digester := seal.NewDigesterSha512(bytes)

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
		var sl *seal.UnwrappedSeal
		sl, err = digester.Unwrap(in, ioutil.Discard)
		fmt.Fprintf(out, "claim:  %v\nactual: %v\n",
			hex.EncodeToString(sl.Claim),
			hex.EncodeToString(sl.Actual))

	case Dump:
		err = seal.DumpHeader(in, out)
	default:
		panic("no command specified")
	}

	return err
}
