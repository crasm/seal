// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

import "errors"

type Command int

const (
	Wrap Command = iota
	Unwrap
	Check
	Dump
)

func getCommand() (Command, error) {
	var cmd Command

	if !isMutuallyExclusive(opt.Wrap, opt.Unwrap, opt.Check, opt.Dump) {
		return cmd, errors.New("too many primary commands")
	}

	switch {
	case opt.Wrap:
		cmd = Wrap
	case opt.Unwrap:
		cmd = Unwrap
	case opt.Check:
		cmd = Check
	case opt.Dump:
		cmd = Dump
	default:
		return cmd, errors.New("no command specified")
	}

	return cmd, nil
}
