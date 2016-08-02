package main

import (
	"os"
	"testing"
)

var stdin = os.Stdin.Name()
var stdout = os.Stdout.Name()

func TestDetermineInputOutput(t *testing.T) {
	cases := []struct {
		givenIn, givenOut       string
		expectedIn, expectedOut string
	}{
		{
			// Explicit os.Stdin and os.Stdout.
			givenIn:     "-",
			givenOut:    "-",
			expectedIn:  stdin,
			expectedOut: stdout,
		},
		{
			// TODO: Are these permutations really necessary?
			// Implicit os.Stdin and os.Stdout.
			givenIn:     "",
			givenOut:    "",
			expectedIn:  stdin,
			expectedOut: stdout,
		},
		{
			// Explicit os.Stdin, implicit os.Stdout.
			givenIn:     "-",
			givenOut:    "",
			expectedIn:  stdin,
			expectedOut: stdout,
		},
		{
			// Implicit os.Stdin, explicit os.Stdout.
			givenIn:     "",
			givenOut:    "-",
			expectedIn:  stdin,
			expectedOut: stdout,
		},
		{
			// Explicit input and output.
			givenIn:     "fileIn",
			givenOut:    "fileOut",
			expectedIn:  "fileIn",
			expectedOut: "fileOut",
		},
		/*
			{
				// Explicit input, implicit output.
				// TODO: How do we know if we need to append or remove
				// ".shd"?
				givenIn:     "fileIn",
				givenOut:    "",
				expectedIn:  "fileIn",
				expectedOut: "fileIn" + ".shd",
			},
		*/
	}

	for _, c := range cases {
		in, out := determineInputOutput(c.givenIn, c.givenOut)
		if in != c.expectedIn {
			t.Fatalf("expected %s, got %s", c.expectedIn, in)
		} else if out != c.expectedOut {
			t.Fatalf("expected %s, got %s", c.expectedOut, out)
		}
	}
}
