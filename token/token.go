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
	MOD   = "%" // todo
	POWER = "^" // todo

	// complex operations
	COS = "cos"
	SIN = "sin"
	// TAN = "tan" // todo / impossible?
)

// reversed keywords
var keywords = map[string]TokenType{
	"cos": COS,
	"sin": SIN,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return ERROR
}
