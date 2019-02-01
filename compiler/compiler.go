package compiler

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

// Compiler holds our object-state.
type Compiler struct {

	// expression holds the mathematical expression we're compiling.
	expression string

	// tokens holds the expression, broken down into a series of tokens.
	//
	// The tokens are received from the lexer, and are not modified.
	tokens []token.Token
}

// New creates a new compiler, given the expression in the constructor.
func New(input string) *Compiler {
	c := &Compiler{expression: input}
	return c
}

// Compiler populates our internal list of tokens, as a result of
// lexing the input string.
//
// There is some error-handling to ensure that the program looks
// somewhat reasonable.
func (c *Compiler) Compile() error {

	//
	// Create the lexer, which will parse our expression.
	//
	lexed := lexer.New(c.expression)

	//
	// First of all populate that `program` array with our tokens.
	//
	for {
		tok := lexed.NextToken()
		if tok.Type == token.EOF {
			break
		}
		if tok.Type == token.ERROR {
			return (fmt.Errorf("Error parsing input; token.ERROR returned from the lexer"))
		}
		c.tokens = append(c.tokens, tok)
	}

	//
	// If the program is empty that's an error.
	//
	if len(c.tokens) < 1 {
		return (fmt.Errorf("The input expression was empty"))
	}

	//
	// If the first token isn't a number we're in trouble
	//
	if c.tokens[0].Type != token.INT {
		return (fmt.Errorf("We expected the program to begin with an integer!\n"))
	}

	//
	// The program should have:
	//   starting-value
	//   NUMBER OPERATOR..
	//   NUMBER OPERATOR..
	//   ..
	//   NUMBER OPERATOR..
	//
	// That means the program should always have an odd-length
	if len(c.tokens)%2 == 0 {
		return (fmt.Errorf("The program should always have an odd-length"))
	}

	//
	// No error.
	//
	return nil
}

// Output converts our series of tokens (i.e. the lexed expression) into
// an assembly-language program.
func (c *Compiler) Output() (string, error) {

	//
	// Get the starting integer.
	//
	start := c.tokens[0].Literal

	//
	// We expect to work in pairs; reading two elements from
	// our program.
	//
	// For example with the program:
	//
	//  3 4 + 5 * 2 /
	//
	// We're going to iterate over the program - skipping the first
	// token - so we'll expect:
	//
	//  4 +
	//  5 *
	//  2 /
	//
	// i.e. Number, operator
	//
	// We'll populate an array of the operations we emit
	// to assembly language.
	//
	var operations []string

	//
	// Temporary-storage for the number we're operating upon next.
	//
	var i string

	//
	// Walk over the array of tokens - skipping the first one,
	// which is our initial value.
	//
	for offset, ent := range c.tokens[1:] {

		//
		// If we have no number then save it.
		//
		if i == "" {
			// number
			i = ent.Literal
			continue
		}

		//
		// The number is already set, so we're now expecting
		// an operator which will be applied to it.
		//
		switch ent.Type {

		case token.PLUS:
			operations = append(operations, `add rax, `+i)

		case token.MOD:

			//
			// Modulus is a cheat - div/idiv will handle
			// setting the remainder in `edx`.  But you
			// need to clear it first to avoid bogus
			// values.
			//
			operations = append(operations, `xor rdx, rdx`)
			operations = append(operations, `mov rax, `+i)
			operations = append(operations, `cqo`)
			operations = append(operations, `div rbx`)
			operations = append(operations, `mov eax, edx`)

		case token.MINUS:
			operations = append(operations, `sub rax,`+i)

		case token.POWER:

			//
			// N ^ 0 -> 0
			// N ^ 1 -> N
			// N ^ 2 -> N * N
			// N ^ 3 -> N * N * N
			// ..
			//
			switch i {
			case "0":
				operations = append(operations, `xor rax, rax`)
			case "1":
				// nop
			default:
				//
				// We'll want to output a loop which
				// means we need a unique label
				//
				// We generate the label ID as the offset of
				// the statement we're generating in our input
				//
				label := fmt.Sprintf("label_%d", offset)

				//
				// Output the loop to calculate the power.
				//
				operations = append(operations, `mov rcx, `+i)
				operations = append(operations, `mov ebx, eax`)
				operations = append(operations, `dec rcx`)
				operations = append(operations, label+":")
				operations = append(operations, `  mul ebx`)
				operations = append(operations, `  dec rcx`)
				operations = append(operations, `  jnz `+label)
			}
		case token.SLASH:
			// Handle a division by zero at run-time.
			// We could catch it at generation-time, just as well..
			if i == "0" {
				operations = append(operations, `jmp div_by_zero`)
			} else {
				operations = append(operations, `mov ebx, `+i)
				operations = append(operations, `cqo`)
				operations = append(operations, `div ebx`)
			}

		case token.ASTERISK:
			operations = append(operations, `mov ebx, `+i)
			operations = append(operations, `mul ebx`)

		default:
			return "", fmt.Errorf("Invalid program - expected operator, but found %v\n", ent)
		}

		//
		// Next time around the loop we'll be looking for
		// a number, rather than an operator.
		//
		i = ""
	}

	//
	// Now we have our starting number, and our list of operations
	//
	// Create a structure to hold these such that we can populate
	// our output-template.
	//
	type Assembly struct {
		Start      string
		Operations []string
	}

	//
	// Create an instance of the output-structure, and populate it.
	//
	var out Assembly

	//
	// Starting value.
	//
	out.Start = start

	//
	// Assembly-language operations
	//
	out.Operations = operations

	//
	// This is the template we'll output.
	//
	assembly := `.intel_syntax noprefix
.global main

.data
format: .asciz "Division by zero\n"
result: .asciz "Result %d\n"

main:
 mov rax, {{.Start}}
{{range .Operations}} {{.}}
{{end}}
 lea rdi,result
 mov rsi, rax
 xor rax, rax
 call printf
 xor rax, rax
 ret

div_by_zero:
 push rbx
 lea  rdi,format
 call printf
 pop rbx
 mov rax, 0
 ret
`

	//
	// Compile the template.
	//
	t := template.Must(template.New("tmpl").Parse(assembly))

	//
	// And now execute it, into a buffer.
	//
	buf := &bytes.Buffer{}
	err := t.Execute(buf, out)
	if err != nil {
		return "", err
	}

	return buf.String(), nil

}
