package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const Stdio = `-`

var opt struct {
	Create  bool `short:"C" long:"create" description:"Create a shielded file."`
	Extract bool `short:"X" long:"extract" description:"Extract a contained file."`
	Verify  bool `short:"V" long:"verify" description:"Verify and check for corruption."`
	Dump    bool `short:"D" long:"dump" description:"Dump raw shield header."`
	// Info    bool `short:"I" long:"info" description:"View shield header information."`

	Output string `short:"o" long:"output" description:"Write output to a file."`
	Force  bool   `short:"f" long:"force" description:"Overwrite files."`
	//Timid      bool `short:"t" long:"timid" description:"Do not allow invalid files to be extracted."`
	//Lax   bool `short:"l" long:"lax" description:"Allow partial and unverified extraction"`
	//Quiet bool `short:"q" long:"quiet" description:"Silence all non-data output to stdout or stderr."`

	Size int `short:"s" long:"size" description:"Truncated size of SHA512 hash in bits." default:"256"`
}

// Slightly complex exit-on-error function. Can handle arbitrary inputs,
// but if the first argument is a string, the remaining arguments can be
// inserted into the string printf-style.
func die(a ...interface{}) {
	if a == nil || len(a) == 0 {
		os.Exit(1)
	}

	buf := bytes.NewBufferString("Error: ")

	switch t := a[0].(type) {
	case string:
		format := t + "\n"
		if len(a) == 1 {
			buf.WriteString(format)
		} else {
			fmt.Fprintf(buf, format, a[1:]...)
		}
	default:
		fmt.Fprintln(buf, a...)
	}

	buf.WriteTo(os.Stderr)
	os.Exit(1)
}

func help(p *flags.Parser) {
	p.WriteHelp(os.Stderr)
	os.Stderr.WriteString("\n")
}

func main() {
	parser := flags.NewParser(&opt, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		help(parser)
		die()
	}

	// Running with no arguments prints help.
	if len(os.Args) == 1 {
		help(parser)
		os.Exit(0)
	}

	in := Stdio
	out := opt.Output

	if len(args) == 1 {
		in = args[0] // We were given an explicit input, so use it. Might still be Stdio.
	} else if len(args) > 1 {
		die("Too many input arguments. Expected only one.")
	}

	if !isMutuallyExclusive(opt.Create, opt.Extract, opt.Verify, opt.Dump) {
		die("You can specify only one primary command at a time.")
	}

	err = dispatch(in, out)
	if err != nil {
		die(err)
	}
}
