package compiler

import "testing"

// TestEscape tests excaping numbers to constants
func TestEscape(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{"3", "const_3"},
		{"0.03", "const_0_03"},
		{"-3", "const_neg_3"},
		{"-3.3", "const_neg_3_3"},
	}

	for _, text := range tests {

		c := New("")

		got := c.escapeConstant(text.input)

		if got != text.expected {
			t.Errorf("Expected '%s' to become '%s', got '%s'",
				text.input, text.expected, got)
		}
	}
}

// TestGenerators just calls the various generating methods, to ensure
// they're covered.
// Since there is no logic in them testing them is pretty pointless.
func TestGenerators(t *testing.T) {

	// create
	c := New("2 3+")

	// misc
	c.genPush("3.4")

	// simple
	c.genPlus()
	c.genMinus()
	c.genMultiply()
	c.genDivide()

	// misc
	c.genModulus()
	c.genPower(1)

	// complex
	c.genAbs()
	c.genCos()
	c.genSin()
	c.genSqrt()
	c.genTan()

	// stack
	c.genDup()
	c.genSwap()
}
