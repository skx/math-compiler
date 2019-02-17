// token package contains the tokens that the lexer will produce when
// parsing an input-expression.
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
	MOD   = "%"
	POWER = "^"

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
var keywords = map[string]TokenType{
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
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return ERROR
}
