// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

import (
	"errors"
	"os"
	"strings"
)

const FileExtension = `.shd`

var stdin = os.Stdin.Name()
var stdout = os.Stdout.Name()

func isExplicit(file string) bool {
	switch file {
	case "", "-", stdin, stdout:
		return false
	default:
		return true
	}
}

// Determine the proper input and output files based on the command and
// file arguments given by the user.
func determineInputOutput(cmd Command, inArg, outArg string) (in, out string, err error) {
	in = stdin
	out = stdout

	infer := false

	if isExplicit(inArg) {
		in = inArg
		infer = outArg == ""
	}

	if isExplicit(outArg) {
		out = outArg
	}

	if infer {
		switch cmd {
		case Create:
			out = in + FileExtension
		case Extract:
			inferred := strings.TrimSuffix(in, FileExtension)
			if inferred == in {
				err = errors.New("output filename required")
			}
			out = inferred
		default:
			// If it's none of the above, leave it as Stdio.
		}
	}

	return in, out, err
}
