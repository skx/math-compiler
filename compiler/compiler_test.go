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
		"3 3",

		// Again
		"3 4 + 3",
	}

	for _, test := range tests {
		c := New(test)
		err := c.tokenize()
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
		"3 sin",
		"4 cos",
		"5 tan",
		"10 sqrt",
		"10 dup +",
		"10 3 swap -",
	}

	for _, test := range tests {

		c := New(test)

		// tokenize
		err := c.tokenize()
		if err != nil {
			t.Errorf("We didn't expect an error tokenizing a valid program, but found one %s", err.Error())
		}

		// convert to internal form
		c.makeinternalform()

		// output the text
		_ = c.output()
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
		err := c.tokenize()
		if err != nil {
			t.Errorf("We didn't expect an error compiling a valid program, but found one %s", err.Error())
		}

		// output
		out := c.output()

		// sanity-check
		if !strings.Contains(out, "main") {
			t.Errorf("Our generated program looked .. bogus?")
		}
	}
}

// Test that enabling the debug-flag generates a trap instruction.
func TestDebug(t *testing.T) {

	test := `1 1 +`

	c := New(test)
	c.SetDebug(true)

	// tokenize
	err := c.tokenize()
	if err != nil {
		t.Errorf("We didn't expect an error tokenizing a valid program, but found one %s", err.Error())
	}

	// convert to internal form
	c.makeinternalform()

	// output the text
	out := c.output()
	if err != nil {
		t.Errorf("We didn't expect an error generating our assembly %s", err.Error())
	}

	if !strings.Contains(out, "int 03") {
		t.Errorf("Debug trap not found!")
	}
}
