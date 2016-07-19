package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import "os"

const DefaultPerm = 0644

// True if at most one is true. All can be false.
func isMutuallyExclusive(bools ...bool) bool {
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
func createFile(path string, force bool) (*os.File, error) {
	callopt := os.O_CREATE | os.O_RDWR
	if force {
		callopt |= os.O_TRUNC
	} else {
		callopt |= os.O_EXCL
	}

	return os.OpenFile(path, callopt, DefaultPerm)
}
