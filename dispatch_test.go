package main

import (
	"os"
	"testing"
)

var stdin = os.Stdin.Name()
var stdout = os.Stdout.Name()

func TestDetermineInputOutput(t *testing.T) {
	t.Parallel()

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
		{
			// Implicit input, explicit output.
			givenIn:     "",
			givenOut:    "fileOut",
			expectedIn:  stdin,
			expectedOut: "fileOut",
		},
	}

	for _, c := range cases {
		for cmd := Create; cmd <= Verify; cmd++ {
			in, out, err := determineInputOutput(cmd, c.givenIn, c.givenOut)
			if err != nil {
				t.Fatalf("expected nil, got %e", err)
			} else if in != c.expectedIn {
				t.Fatalf("expected %s, got %s", c.expectedIn, in)
			} else if out != c.expectedOut {
				t.Fatalf("expected %s, got %s", c.expectedOut, out)
			}
		}
	}
}

func TestDetermineInputOutputInference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command                 Command
		givenIn, givenOut       string
		expectedIn, expectedOut string
	}{
		{
			command:     Create,
			givenIn:     "fileIn",
			givenOut:    "",
			expectedIn:  "fileIn",
			expectedOut: "fileIn.shd",
		},
		{
			command:     Extract,
			givenIn:     "fileIn.shd",
			givenOut:    "",
			expectedIn:  "fileIn.shd",
			expectedOut: "fileIn",
		},
	}

	for _, c := range cases {
		in, out, err := determineInputOutput(c.command, c.givenIn, c.givenOut)
		if err != nil {
			t.Fatalf("expected nil, got %e", err)
		} else if in != c.expectedIn {
			t.Fatalf("expected %s, got %s", c.expectedIn, in)
		} else if out != c.expectedOut {
			t.Fatalf("expected %s, got %s", c.expectedOut, out)
		}
	}

}

func TestDetermineInputOutputExtractMissingExtension(t *testing.T) {
	t.Parallel()

	in, out, err := determineInputOutput(Extract, "fileIn", "")
	if err == nil {
		t.Fatalf("expected an error, got in = '%s', out='%s'", in, out)
	}
}
