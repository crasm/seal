// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package seal

import (
	"errors"
	"io"
	"os"
)

var ErrFileModified = errors.New("seal: file has been modified without a Sync")

// File represents a subset of the functionality in os.File, to allow
// transparent access to reading and writing seal files.
type File interface {
	io.Reader
	// io.ReaderAt
	io.Writer
	// io.WriterAt
	io.Closer
	io.Seeker
	// Returns
	Validate() error
	// Name() string
	// Sync() error
	// Truncate (size int64) error
}

type file struct {
	osFile os.File
}
