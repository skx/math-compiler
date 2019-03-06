// Package token contains the tokens that the lexer will produce when
// parsing an input-expression.
package token

// Type is a string
type Type string

// Token struct represent the lexer token
type Token struct {
	Type    Type
	Literal string
}

// pre-defined Type
const (
	EOF    = "EOF"
	ERROR  = "ERROR"
	NUMBER = "NUMBER"
	IDENT  = "IDENT"

	// simple operations
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	// advanced operations
	MOD       = "%"
	POWER     = "^"
	FACTORIAL = "!"

	// misc
	E  = "e"
	PI = "pi"

	// complex operations
	ABS  = "abs"
	COS  = "cos"
	SIN  = "sin"
	SQRT = "sqrt"
	TAN  = "tan"

	// stack operations
	DUP  = "dup"
	SWAP = "swap"
)

// reversed keywords
var keywords = map[string]Type{
	"abs":  ABS,
	"cos":  COS,
	"dup":  DUP,
	"e":    E,
	"pi":   PI,
	"sin":  SIN,
	"sqrt": SQRT,
	"swap": SWAP,
	"tan":  TAN,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) Type {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return ERROR
}
