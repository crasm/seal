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
	var err error

	switch cmd {
	case Wrap:
		if out.Name() == os.Stdout.Name() {
			_, err = seal.WrapBufferedBits(in, out, opt.Size)
		} else {
			_, err = seal.WrapBits(in, out, opt.Size)
		}

	case Unwrap:
		_, err = seal.Unwrap(in, out)

	case Check:
		var sl *seal.UnwrappedSeal
		sl, err = seal.Unwrap(in, ioutil.Discard)
		fmt.Fprintf(out, "claim:  %v\nactual: %v\n",
			hex.EncodeToString(sl.ClaimedSignature),
			hex.EncodeToString(sl.CalculatedSignature))

	case Dump:
		err = seal.DumpHeader(in, out)
	default:
		panic("no command specified")
	}

	return err
}
