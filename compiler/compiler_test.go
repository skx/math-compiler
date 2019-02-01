package compiler

import (
	"testing"
)

// We try to compile several bogus programs
func TestBogusInput(t *testing.T) {

	tests := []string{

		// empty program
		"",

		// program that doesn't start with an int
		"+",

		// program with invalid token
		"3 5 $",

		// program with a missing operator
		"3 3",

		// Again
		"3 4 + 3",
	}

	for _, test := range tests {
		c := New(test)
		err := c.Compile()
		if err == nil {
			t.Errorf("We expected an error handling '%s', but got none!", test)
		}
	}
}
