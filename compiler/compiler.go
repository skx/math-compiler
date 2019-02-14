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
	"fmt"
	"strconv"
	"strings"

	"github.com/skx/math-compiler/instructions"
	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

// Compiler holds our object-state.
type Compiler struct {

	// expression holds the mathematical expression we're compiling.
	expression string

	//
	// Constants we come across.
	//
	// Rather than placing the numbers in-line we're storing them
	// in the constant area.  This gives us a level of indirection,
	// but it simpler to reason about.
	//
	// Note that if we see "3 + 4 + 3 + 4" we only store each of
	// the constants once each.  So we have a map here, key is the
	// value of the constant, and `bool` to record that we found it.
	//
	// Later we'll output just the unique keys.
	//
	constants map[string]bool

	// tokens holds the expression, broken down into a series of tokens.
	//
	// The tokens are received from the lexer, and are not modified.
	tokens []token.Token

	// Instructions is the virtual instructions we're going to compile
	// to assembly
	instructions []instructions.Instruction
}

// New creates a new compiler, given the expression in the constructor.
func New(input string) *Compiler {
	c := &Compiler{expression: input, constants: make(map[string]bool)}
	return c
}

// Tokenize populates our internal list of tokens, as a result of
// lexing the input string.
//
// There is some error-handling to ensure that the program looks
// somewhat reasonable.
func (c *Compiler) Tokenize() error {

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
	// Get the last token
	//
	if len(c.tokens) > 1 {
		len := len(c.tokens)
		end := c.tokens[len-1]
		if end.Type == token.NUMBER {
			return fmt.Errorf("Program ends with a number, which is invalid!")
		}
	}

	//
	// No error.
	//
	return nil
}

// InternalForm converts our series of tokens (i.e. the lexed expression) into
// an an intermediary form, collecting constants as we go.
func (c *Compiler) InternalForm() error {

	//
	// Walk our tokens.
	//
	for _, t := range c.tokens {

		switch t.Type {
		case token.NUMBER:

			// Mark the constant as having been used.
			// Record that this constant is set.
			c.constants[t.Literal] = true

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Push, Value: t.Literal})

		case token.PLUS:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Plus, Args: 2})

		case token.MOD:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Modulus, Args: 2})

		case token.MINUS:

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Minus, Args: 2})

		case token.POWER:

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Power, Args: 2})

		case token.SLASH:

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Divide, Args: 2})

		case token.ASTERISK:

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Multiply, Args: 2})

		case token.SIN:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Sin, Args: 1})

		case token.TAN:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Tan, Args: 1})

		case token.SQRT:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Sqrt, Args: 1})

		case token.COS:
			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Instruction: instructions.Cos, Args: 1})

		default:
			return fmt.Errorf("Invalid program - expected operator, but found %v\n", t)
		}
	}

	return nil

}

// Output writes our program to stdout
func (c *Compiler) Output() (string, error) {

	//
	// The header.
	//
	header := `
.intel_syntax noprefix
.global main

# Data-section: Contains the format-string for our output message,
#               etc.
.data
   int: .double 0.0
   a: .double 0.0
   b: .double 0.0
   fmt:   .asciz "Result %g\n"
`

	//
	// Add on the constants
	//
	for v, _ := range c.constants {
		header += fmt.Sprintf("%s: .double %s\n",
			c.escapeConstant(v), v)
	}

	header += `main:
push rbp
`

	//
	// The body of the program
	//
	body := ""

	//
	// Now we walk over our form.
	//
	for _, opr := range c.instructions {
		switch opr.Instruction {

		case instructions.Push:
			body += fmt.Sprintf(`
        fld qword ptr %s
        fstp qword ptr [int]
        mov rax, qword ptr [int]
        push rax
`, c.escapeConstant(opr.Value))

		case instructions.Plus:
			body += c.genPlus()

		case instructions.Minus:
			body += c.genMinus()

		case instructions.Multiply:
			body += c.genMultiply()

		case instructions.Divide:
			body += c.genDivide()

		case instructions.Power:
			body += c.genNop()

		case instructions.Modulus:
			body += c.genNop()

		case instructions.Sin:
			body += c.genNop()

		case instructions.Cos:
			body += c.genNop()

		case instructions.Tan:
			body += c.genNop()

		case instructions.Sqrt:
			body += c.genNop()

		}
	}

	footer := `
        # print the result
        lea rdi,fmt

        # get the value to print in xmm0
        pop rax
        mov qword ptr [a], rax
        movq xmm0, [a]

        movq rax, 1


        call printf

        # clean and exit
        pop	rbp
        xor rax, rax
        ret
`

	return header + body + footer, nil
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

	var val string
	if s < 0 {
		val = fmt.Sprintf("const_neg_%s", input)
	} else {
		val = fmt.Sprintf("const_%s", input)
	}

	val = strings.Replace(val, ".", "_", -1)
	return val
}

func (c *Compiler) genPlus() string {

	return `
        # pop two values
        pop rax
        mov qword ptr [a], rax
        pop rax
        mov qword ptr [b], rax

        # add
        fld qword ptr [a]
        fadd qword ptr  [b]
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`

}
func (c *Compiler) genMinus() string {
	return `
        # pop two values
        pop rax
        mov qword ptr [a], rax
        pop rax
        mov qword ptr [b], rax

        # sub
        fld qword ptr [b]
        fsub qword ptr  [a]
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`
}
func (c *Compiler) genMultiply() string {
	return `
        # pop two values
        pop rax
        mov qword ptr [a], rax
        pop rax
        mov qword ptr [b], rax

        # multiply
        fld qword ptr [a]
        fmul qword ptr  [b]
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`

}
func (c *Compiler) genDivide() string {
	return `
        # pop two values
        pop rax
        mov qword ptr [a], rax
        pop rax
        mov qword ptr [b], rax

        # divide
        fld qword ptr [b]
        fdiv qword ptr  [a]
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`

}

func (c *Compiler) genNop() string {
	return ``
}
