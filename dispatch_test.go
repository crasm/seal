package main

import (
	"os"
	"testing"
)

func TestDetermineInputOutput(t *testing.T) {
	cases := []struct {
		givenInArg      string
		givenOutArg     string
		expectedInFile  *os.File
		expectedOutFile *os.File
	}{
		{
			givenInArg:      "",
			givenOutArg:     "",
			expectedInFile:  os.Stdin,
			expectedOutFile: os.Stdout,
		},
	}

	for _, c := range cases {
		inFile, outFile, err := determineInputOutput(c.givenInArg, c.givenOutArg)
		if err != nil {
			t.Fatal(err)
		} else if inFile != c.expectedInFile {
			t.Fatalf("expected %s, got %s", c.expectedInFile, inFile)
		} else if outFile != c.expectedOutFile {
			t.Fatalf("expected %s, got %s", c.expectedOutFile, outFile)
		}
	}
}
