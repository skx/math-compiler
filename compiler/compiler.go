// compiler contains our simple compiler.
//
// In brief it uses the lexer to tokenize the expression, and outputs
// a simple assembly language program.
//
// There are only two minor complications:
//
//  1.  We store all the input-floats in the data-area of the program.
//      These require escaping for uniqueness purposes.
//
//  2.  We output different instructions based on the operator.
//
//  we could do better with an intermediary form, concatenating small
// segments of code, and avoiding the frequent loads from the `result`
// location.
//
// That said this is a toy, and will remain a toy, so I can live with
// these problems.
//

package compiler

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
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
		// Get the next token.
		tok := lexed.NextToken()

		// If end of stream then break.
		if tok.Type == token.EOF {
			break
		}

		// If error then abort.
		if tok.Type == token.ERROR {
			return (fmt.Errorf("Error parsing input; token.ERROR returned from the lexer"))
		}

		// Otherwise append the token to our program.
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
	if c.tokens[0].Type != token.NUMBER {
		return (fmt.Errorf("We expected the program to begin with a numeric thing!\n"))
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
	// Constants we come across.
	//
	// Rather than placing the integers in-line we're storing them
	// in the constant area.  This gives us a level of indirection,
	// but it simpler to reason about.
	//
	// Note that if we see "3 + 4 + 3 + 4" we only store each of
	// the constants once each.  So we have a map here, key is the
	// value of the constant, and `bool` to record that we found it.
	//
	// Later we'll output just the unique keys.
	//
	constants := make(map[string]bool)

	//
	// Temporary-storage for the number we're operating upon next.
	//
	var i string

	//
	// Walk over the array of tokens - skipping the first one,
	// which is our initial value.
	//
	for _, ent := range c.tokens[1:] {

		switch ent.Type {

		case token.NUMBER:
			// Save the literal value away
			i = ent.Literal

			// Record that this constant is set.
			constants[ent.Literal] = true

		case token.PLUS:

			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// add the constant
			constant := c.escapeConstant(i)
			operations = append(operations, `fadd qword ptr [`+constant+`]`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		case token.MOD:

			return "", fmt.Errorf("token.MOD - not implemented")

		case token.MINUS:

			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// subtract the constant
			constant := c.escapeConstant(i)
			operations = append(operations, `fsub qword ptr [`+constant+`]`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		case token.POWER:
			return "", fmt.Errorf("token.POWER - not implemented")

		case token.SLASH:

			if i == "0" {
				fmt.Printf("Division by zero!")
				os.Exit(1)

			}

			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// divide by the constant
			constant := c.escapeConstant(i)
			operations = append(operations, `fdiv qword ptr [`+constant+`]`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		case token.ASTERISK:

			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// multiply by the constant
			constant := c.escapeConstant(i)
			operations = append(operations, `fmul qword ptr [`+constant+`]`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		case token.SIN:

			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// run the sin
			operations = append(operations, `fsin`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		case token.COS:
			// load the (current) result
			operations = append(operations, `fld qword ptr [result]`)

			// run the cos
			operations = append(operations, `fcos`)

			// store back in the result store.
			operations = append(operations, `fstp qword ptr [result]`)

		default:
			return "", fmt.Errorf("Invalid program - expected operator, but found %v\n", ent)
		}
	}

	//
	// Now we have our starting number, and our list of operations
	//
	// Create a structure to hold these such that we can populate
	// our output-template.
	//
	type Assembly struct {
		// The starting value.
		Start string

		// The operations we carry out.
		Operations []string

		// Any constants we load.
		Constants []string
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
	// The constants
	//
	for key, _ := range constants {
		out.Constants = append(out.Constants, key)
	}

	//
	// This is the template we'll output.
	//
	assembly := `
.intel_syntax noprefix
.global main

# Data-section: Contains the format-string for our output message,
#               the starting value, and any constants we use.
.data
result: .double {{.Start}}
fmt:   .asciz "Result %g\n"
{{range .Constants}}const_{{.}}: .double {{.}}
{{end}}

main:
        push	rbp

{{range .Operations}} {{.}}
{{end}}

        # print the result
        lea rdi,fmt
        movq rax, 1
        movq xmm0, [result]
        call printf

        # clean and exit
        pop	rbp
        xor rax, rax
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

// escapeConstant converts a floating-point number such as
// "1.2", or "-1.3" into a constant value that can be embedded
// safely into our generated assembly-language file.
func (c *Compiler) escapeConstant(input string) string {

	// Convert "3.0" to "const_3.0"
	s, err := strconv.ParseFloat(input, 32)

	if err != nil {
		fmt.Printf("Failed to parse '%s' into a float\n", input)
		fmt.Printf("%s\n", err.Error())
	}

	if s < 0 {
		return fmt.Sprintf("const_neg_%s", input)
	}
	return fmt.Sprintf("const_%s", input)
}
