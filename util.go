// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

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
