package main

import "testing"

func TestIsMutuallyExclusive(t *testing.T) {
	t.Parallel()

	cases := []struct {
		bools  []bool
		answer bool
	}{
		{
			bools:  []bool{},
			answer: true,
		},
		{
			bools:  []bool{true},
			answer: true,
		},
		{
			bools:  []bool{false},
			answer: true,
		},
		// TODO rest of truth table, or is that overkill?
		{
			bools:  []bool{false, false, false},
			answer: true,
		},
		{
			bools:  []bool{false, false, true},
			answer: true,
		},
		{
			bools:  []bool{false, true, true},
			answer: false,
		},
		{
			bools:  []bool{true, true, true},
			answer: false,
		},
	}

	for i, c := range cases {
		if a := isMutuallyExclusive(c.bools...); a != c.answer {
			t.Errorf("case %d: expected %t, got %t", i, c.answer, a)
		}
	}
}
