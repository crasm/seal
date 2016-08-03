// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package shield

import "io"

// A Digester knows how to read and write entire shield files in a
// single pass.
type Digester interface {
	// Wrap the contents of the io.Reader and write it to the
	// io.WriteSeeker. Because of the ability to seek the output, this
	// method takes linear time and space, so should be preferred.
	Wrap(io.Reader, io.WriteSeeker) (*Shield, error)
	// Wrap the contents of the io.Reader, but buffer the output to a
	// temporary file while digesting the input.
	WrapBuffered(io.Reader, io.Writer) (*Shield, error)
	Unwrap(io.Reader, io.Writer) (*UnwrappedShield, error)
}
