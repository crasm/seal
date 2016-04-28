package main

// Copyright (c) 2016, Christian Demsar
// This code is open source under the ISC license. See LICENSE for details.

import (
	"github.com/crasm/shield"
	"github.com/jessevdk/go-flags"
)

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const DefaultPerm = 0644
const FileExtension = `.shd`

const Stdio = `-`

var opt struct {
	Create  bool `short:"C" long:"create" description:"Create a shield file."`
	Extract bool `short:"X" long:"extract" description:"Extract a contained file from a shield file."`
	Info    bool `short:"I" long:"info" description:"Show info on a shield file."`

	Output    string `short:"o" long:"output" description:"Write output to a file."`
	inferName bool

	Force bool `short:"f" long:"force" description:"Overwrite files."`
	//Timid      bool `short:"t" long:"timid" description:"Delete extracted file if its claim is found to be invalid."`
	//Lax   bool `short:"l" long:"lax" description:"Allow partial and unverifieid dextraction"`
	//Quiet bool `short:"q" long:"quiet" description:"Silence all non-data output to stdout or stderr."`
}

// Figures out input and output files and calls the appropiate shield library
// functions on them.
func dispatch(in, out string) error {
	var err error

	inFile := os.Stdin
	outFile := os.Stdout

	if in != Stdio && in != "" {
		inFile, err = os.Open(in)
		defer inFile.Close()
		if err != nil {
			return err
		}
	}

	inferName := out == "" && in != "" && in != Stdio

	switch {
	case opt.Create:
		if inferName {
			inferred := fmt.Sprint(in, FileExtension)
			outFile, err = safeFileCreate(inferred)
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
		if inferName {
			inferred := strings.TrimSuffix(in, FileExtension)
			outFile, err = safeFileCreate(inferred)
			defer outFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
		err = shield.Unwrap(inFile, outFile)
	case opt.Info:
		fallthrough
	default:
		panic("info command not supported (yet)")
	}

	return err
}

func main() {
	args, err := flags.Parse(&opt)
	if err != nil {
		log.Fatal(err)
	}

	if !xor(opt.Create, opt.Extract, opt.Info) {
		log.Fatal("more than one command (or no commands) specified")
	}

	if len(args) > 1 {
		log.Fatal("can work on at most a single shield file")
	}

	in := Stdio
	out := opt.Output

	if len(args) == 1 { // If given an input file, use that. Might still be Stdio.
		in = args[0]
	} else { // No input or output files. Assume Stdio for both input and output.
		out = Stdio
	}

	err = dispatch(in, out)
	if err != nil {
		log.Fatal(err)
	}
}

// True if at most one is true.
func xor(bools ...bool) bool {
	found := 0
	for _, b := range bools {
		if b {
			found++
		}
	}
	return found <= 1
}

// Creates the file at the given location. If opt.Force is set, the existing
// file is clobbered.
func safeFileCreate(path string) (*os.File, error) {
	callopt := os.O_CREATE | os.O_RDWR
	if opt.Force {
		callopt |= os.O_TRUNC
	} else {
		callopt |= os.O_EXCL
	}

	return os.OpenFile(path, callopt, DefaultPerm)
}
