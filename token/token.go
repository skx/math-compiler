package token

// TokenType is a string
type TokenType string

// Token struct represent the lexer token
type Token struct {
	Type    TokenType
	Literal string
}

// pre-defined TokenType
const (
	EOF      = "EOF"
	ERROR    = "ERROR"
	INT      = "INT"
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"
	POWER    = "^"
)
