package main

// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

import (
	"fmt"
	"os"
)

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

// Creates the file at the given location. If force is set, the existing
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

func bitsToBytes(bits int) (int, error) {
	if bits < 0 {
		return -1, fmt.Errorf("bits was negative: %v", bits)
	}

	bytes := bits / 8
	if bytes*8 != bits {
		return -1, fmt.Errorf("bits was not a multiple of 8: %v", bits)
	}

	return bytes, nil
}
