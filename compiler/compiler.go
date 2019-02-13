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

	"github.com/skx/math-compiler/instructions"
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

	// Instructions is the virtual instructions we're going to compile
	// to assembly
	instructions []instructions.Instruction
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

// Output converts our series of tokens (i.e. the lexed expression) into
// an assembly-language program.
func (c *Compiler) Output() (string, error) {

	//
	// Walk our tokens.
	//
	for _, t := range c.tokens {

		switch t.Type {
		case token.NUMBER:

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
			return "", fmt.Errorf("Invalid program - expected operator, but found %v\n", t)
		}
	}

	//
	// Show what we've compiled.
	//
	for _, v := range c.instructions {
		fmt.Printf("Type: %c Argument-Count:%d Value:%s\n", v.Instruction, v.Args, v.Value)
	}
	return "", nil

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
