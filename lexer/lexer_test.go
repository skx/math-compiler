package lexer

import (
	"testing"

	"github.com/skx/math-compiler/token"
)

// Trivial test of the parsing of numbers.
func TestParseNumbers(t *testing.T) {
	input := `3 43 -17 -3`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NUMBER, "3"},
		{token.NUMBER, "43"},
		{token.NUMBER, "-17"},
		{token.NUMBER, "-3"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Trivial test of parsing floats.
func TestParseFloats(t *testing.T) {
	input := `3.14 4.3 -1.7 -2.13 sin `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NUMBER, "3.14"},
		{token.NUMBER, "4.3"},
		{token.NUMBER, "-1.7"},
		{token.NUMBER, "-2.13"},
		{token.SIN, "sin"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Trivial test of the parsing of operators.
func TestParseOperators(t *testing.T) {
	input := `+ - * / % ^ -`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.MOD, "%"},
		{token.POWER, "^"},
		{token.MINUS, "-"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Trivial test of the parsing invalid input
func TestParseBogus(t *testing.T) {
	input := `steve 3`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ERROR, "Unknown token steve"},
		{token.NUMBER, "3"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
