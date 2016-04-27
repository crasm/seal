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

var opt struct {
	Create  bool `short:"C" long:"create" description:"Create a shield file."`
	Extract bool `short:"X" long:"extract" description:"Extract a contained file from a shield file."`
	Info    bool `short:"I" long:"info" description:"Show info on a shield file."`

	InferName bool `short:"i" long:"infer-name" description:"Infer output filename."`
	Force     bool `short:"f" long:"force" description:"Overwrite files."`
	//Timid      bool `short:"t" long:"timid" description:"Delete extracted file if its claim is found to be invalid."`
	//Lax   bool `short:"l" long:"lax" description:"Allow partial and unverifieid dextraction"`
	//Quiet bool `short:"q" long:"quiet" description:"Silence all non-data output to stdout or stderr."`
}

func main() {
	args, err := flags.Parse(&opt)
	if err != nil {
		log.Fatal(err)
	}

	if !xor(opt.Create, opt.Extract, opt.Info) {
		log.Fatal("more than one command (or no commands) specified")
	}

	in := os.Stdin
	out := os.Stdout

	if opt.InferName {
		if len(args) != 1 {
			log.Fatal("can only work on a single shield file")
		}
		in, err = os.Open(args[0])
		defer in.Close()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		if len(args) > 2 {
			log.Fatal("more than two non-flag args remaining")
		}

		if len(args) == 2 && args[1] != "-" {
			out, err = safeFileCreate(args[1])
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(args) >= 1 && args[0] != "-" {
			in, err = os.Open(args[0])
			defer in.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	switch {
	case opt.Create:
		if opt.InferName {
			out, err = safeFileCreate(fmt.Sprint(in.Name(), FileExtension))
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
		if out == os.Stdout {
			err = shield.WrapBuffered(in, out)
		} else {
			err = shield.Wrap(in, out)
		}
	case opt.Extract:
		if opt.InferName {
			out, err = safeFileCreate(strings.TrimSuffix(in.Name(), FileExtension))
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
		err = shield.Unwrap(in, out)
	case opt.Info:
		fallthrough
	default:
		panic("info command not supported (yet)")
	}

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
