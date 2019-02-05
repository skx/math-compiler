package compiler

import (
	"strings"
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
		//TODO		"3 3",

		// Again
		//TODO TODO	"3 4 + 3",

	}

	for _, test := range tests {
		c := New(test)
		err := c.Compile()
		if err == nil {
			t.Errorf("We expected an error handling '%s', but got none!", test)
		}
	}
}

// Test some valid programs.
func TestValidPrograms(t *testing.T) {

	tests := []string{
		"1 2 -",
		"3 4 +",
		"5 7 *",
		"9 3 /",
		"10 5 %",
		"2 8 ^",
	}

	for _, test := range tests {
		c := New(test)
		err := c.Compile()
		if err != nil {
			t.Errorf("We didn't expect an error compiling a vlid program, but found one %s", err.Error())
		}
	}
}

// Test actually outputing some valid programs.
//
// This test covers the full range:
//   "parse".
//   "compile".
//   "output".
//
// However it doesn't test that the generated output contains what we
// expect.  The only way to do that would be to have a static-file and
// compare it literally.  If we did that we'd have a pain keeping it
// in sync.
//
// So here we're just looking for rough-behaviour.  Sorry!
//
func TestValidOutput(t *testing.T) {

	tests := []string{
		"1 2 -",
		"3 4 +",
		"5 7 *",
		"9 3 /",
		"10 5 %",
		"2 8 ^",
		"2 0 ^",  // N ^ 0 is a special case
		"2 1 ^",  // N ^ 1 is a special case
		"2 12 ^", // N ^ 12 is NOT a special case!
	}

	for _, test := range tests {

		// create
		c := New(test)

		// compile
		err := c.Compile()
		if err != nil {
			t.Errorf("We didn't expect an error compiling a valid program, but found one %s", err.Error())
		}

		// output
		out := ""
		out, err = c.Output()
		if err != nil {
			t.Errorf("We didn't expect an error outputing a valid program, but found one %s", err.Error())
		}

		// sanity-check
		if !strings.Contains(out, "main") {
			t.Errorf("Our generated program looked .. bogus?")
		}
	}
}

// Test actually outputing some invalid programs.
//
// This test covers the full range:
//   "parse".
//   "compile".
//   "output".
//
func TestInvalidOutput(t *testing.T) {

	tests := []string{
		"3 0 /",
		"3 3.3 %",
		"2 3.4 ^",
	}

	for _, test := range tests {

		// create
		c := New(test)

		// compile
		err := c.Compile()
		if err != nil {
			t.Errorf("We didn't expect an error compiling an invalid program, but found one %s", err.Error())
		}

		// output
		_, err = c.Output()
		if err == nil {
			t.Errorf("We expected an error outputing an invalid program, but found none")
		}
	}
}
