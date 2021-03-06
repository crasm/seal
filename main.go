// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var opt struct {
	Wrap   bool `short:"W" long:"wrap" description:"Wrap a file in a seal."`
	Unwrap bool `short:"U" long:"unwrap" description:"Unwrap (extract) a sealed file."`
	Check  bool `short:"C" long:"check" description:"Check a seal for corrupted file contents."`
	Dump   bool `short:"D" long:"dump" description:"Dump raw seal header."`
	// Info    bool `short:"I" long:"info" description:"View seal header information."`

	Output  string `short:"o" long:"output" description:"Write output to a file."`
	Verbose bool   `short:"v" long:"verbose" description:"Enable verbose debug output"`

	Force bool `long:"force" description:"Overwrite files. Required when inferring filenames."`

	//Timid      bool `short:"t" long:"timid" description:"Do not allow invalid files to be extracted."`
	//Lax   bool `short:"l" long:"lax" description:"Allow partial and unverified extraction"`
	//Quiet bool `short:"q" long:"quiet" description:"Silence all non-data output to stdout or stderr."`

	Size int `short:"s" long:"size" description:"Truncated size of SHA512 hash in bits." default:"256"`

	Debug bool `long:"debug" description:"Log debug information."`
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
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type != flags.ErrHelp {
			help(parser)
		}
		die()
	}

	// Running with no arguments prints help.
	if len(os.Args) == 1 {
		help(parser)
		os.Exit(0)
	}

	// Figure out what we're supposed to do.
	cmd, err := getCommand()
	if err != nil {
		die(err)
	}

	inArg := ""
	outArg := opt.Output

	// If we explicitly give an output file, don't stop us from
	// overwriting it.
	if outArg != "" {
		opt.Force = true
	}

	if len(args) == 1 {
		// We were given an explicit input, so use it. Might still be stdio ("-").
		inArg = args[0]
	} else if len(args) > 1 {
		die("Too many input arguments. Expected only one.")
	}

	in, out, err := determineInputOutput(cmd, inArg, outArg)
	if err != nil {
		die(err)
	}

	if opt.Verbose {
		log.Printf("Using %q for input, %q for output\n", in, out)
	}

	inFile, outFile, err := openInputOutput(cmd, opt.Force, in, out)
	defer inFile.Close()
	defer outFile.Close()

	if err != nil {
		die(err)
	}

	err = dispatch(cmd, inFile, outFile)
	if err != nil {
		die(err)
	}
}
