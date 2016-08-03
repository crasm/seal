package main

import "os"

const DefaultPerm = 0644

func openInputOutput(cmd Command, force bool, in, out string) (inFile, outFile *os.File, err error) {
	inFile, err = os.Open(in)
	if err != nil {
		return
	}

	if out == os.Stdout.Name() {
		outFile, err = os.OpenFile(out, os.O_WRONLY|os.O_APPEND, DefaultPerm)
		return
	}

	// If we got here, we're actually creating a new file!

	callopt := os.O_CREATE | os.O_RDWR
	if force {
		callopt |= os.O_TRUNC
	} else {
		callopt |= os.O_EXCL
	}

	outFile, err = os.OpenFile(out, callopt, DefaultPerm)

	return
}
