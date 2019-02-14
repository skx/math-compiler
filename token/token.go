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

	// complex operations
	COS  = "cos"
	SIN  = "sin"
	SQRT = "sqrt"
	TAN  = "tan"
)

// reversed keywords
var keywords = map[string]TokenType{
	"cos":  COS,
	"sin":  SIN,
	"sqrt": SQRT,
	"tan":  TAN,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return ERROR
}
