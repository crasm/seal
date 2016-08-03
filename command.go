// Copyright (c) 2016, crasm <crasm@vczf.io>
// This code is open source under the ISC license. See LICENSE for details.

package main

import "errors"

type Command int

const (
	Create Command = iota
	Extract
	Verify
	Dump
)

func getCommand() (Command, error) {
	var cmd Command

	if !isMutuallyExclusive(opt.Create, opt.Extract, opt.Verify, opt.Dump) {
		return cmd, errors.New("too many primary commands")
	}

	switch {
	case opt.Create:
		cmd = Create
	case opt.Extract:
		cmd = Extract
	case opt.Verify:
		cmd = Verify
	case opt.Dump:
		cmd = Dump
	default:
		return cmd, errors.New("no command specified")
	}

	return cmd, nil
}
