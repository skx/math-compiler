package lexer

import (
	"strings"

	"github.com/skx/math-compiler/token"
)

// Lexer holds our object-state.
type Lexer struct {
	position     int    //current character position
	readPosition int    //next character position
	ch           rune   //current character
	characters   []rune //rune slice of input string
}

// New a Lexer instance from string input.
func New(input string) *Lexer {
	l := &Lexer{characters: []rune(input)}
	l.readChar()
	return l
}

// read one forward character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case rune('+'):
		tok = newToken(token.PLUS, l.ch)
	case rune('%'):
		tok = newToken(token.MOD, l.ch)
	case rune('^'):
		tok = newToken(token.POWER, l.ch)
	case rune('-'):
		// "-3" is "-3", "-3.4" is "-3.4", but "3 - 4" is -1 (via the distinct tokens "3", "-", "4".)
		if isDigit(l.peekChar()) {

			// swallow the -
			l.readChar()

			// read an int/float
			tok = l.readDecimal()

			// ensure the sign is not lost.
			tok.Literal = "-" + tok.Literal

		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case rune('/'):
		tok = newToken(token.SLASH, l.ch)
	case rune('*'):
		tok = newToken(token.ASTERISK, l.ch)
	case rune(0):
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isDigit(l.ch) {
			return l.readDecimal()
		}

		lit := l.readIdentifier()
		tok.Type = token.LookupIdentifier(lit)
		if tok.Type == token.ERROR {
			tok.Literal = "Unknown token " + lit
		} else {
			tok.Literal = lit
		}
	}
	l.readChar()
	return tok
}

// return new token
func newToken(tokenType token.Type, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// skip white space
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// readNumber handles reading a number, comprising of digits 0-9.
func (l *Lexer) readNumber() string {
	str := ""

	// We only accept digits.
	accept := "0123456789"

	for strings.Contains(accept, string(l.ch)) {
		str += string(l.ch)
		l.readChar()
	}
	return str
}

// read a decimal / floating point number.
func (l *Lexer) readDecimal() token.Token {

	//
	// Read an integer-number.
	//
	integer := l.readNumber()

	//
	// We might have more content:
	//
	//   .[digits]  -> Which converts us from an int to a float.
	//
	if l.ch == rune('.') && isDigit(l.peekChar()) {

		//
		// OK here we think we've got a float.
		//
		// Skip the period.
		//
		l.readChar()

		//
		// Read the fractional part.
		//
		fraction := l.readNumber()
		return token.Token{Type: token.NUMBER, Literal: integer + "." + fraction}
	}
	return token.Token{Type: token.NUMBER, Literal: integer}

}

// peek character
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.readPosition]
}

// is white space
func isWhitespace(ch rune) bool {
	return ch == rune(' ') || ch == rune('\t') || ch == rune('\n') || ch == rune('\r')
}

// is Digit
func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

// readIdentifier is designed to read an identifier which means a string
// such as `sin`, `cos`, `tan`.
func (l *Lexer) readIdentifier() string {

	id := ""

	//
	// Build up our identifier, handling only valid characters.
	//
	for isIdentifier(l.ch) {
		id += string(l.ch)
		l.readChar()
	}

	return id
}

// determinate ch is identifier or not
func isIdentifier(ch rune) bool {
	return !isDigit(ch) && !isWhitespace(ch) && ch != rune(0)
}
