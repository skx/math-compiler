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
		{token.INT, "3"},
		{token.INT, "43"},
		{token.INT, "-17"},
		{token.INT, "-3"},
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
		{token.ERROR, "s"},
		{token.ERROR, "t"},
		{token.ERROR, "e"},
		{token.ERROR, "v"},
		{token.ERROR, "e"},
		{token.INT, "3"},
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
