// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package shield

const Magic = `SHD%`
const Version = 0 // Current major version number for the library.

const IdentLen = len(`SHD%v0`)

// Shield is the information read from a shield header.
type Shield struct {
	Magic   string
	Version int
	Claim   []byte
}

// Shield contains shield file data.
// TODO: Clarify.
type UnwrappedShield struct {
	Shield
	Actual []byte
}
