// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package seal

import "errors"

const Magic = `SL%v`
const Version = 0 // Current major version number for the library.

const IdentLen = len(`SL%v0`)

var ErrSealBroken = errors.New("seal: claim did not validate against content")

// Seal is the information read from a seal header.
type Seal struct {
	Magic   string
	Version int
	Claim   []byte
}

// UnwrappedSeal contains seal file data.
// TODO: Clarify.
type UnwrappedSeal struct {
	Seal
	Actual []byte
}
