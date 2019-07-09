// The compiler-package contains the core of our compiler.
//
// In brief we go through a three-step process:
//
//  1.  It use the lexer to tokenize the expression.
//
//  2.  Convert our program from a series of tokens to an internal form.
//
//  3.  Walk our internal form, generating output for each block.
//
// There are only one minor complication - storing all the input-floats
// in the data-area of the program.  These require escaping for uniqueness
// purposes, and to avoid issues with the assembling.
//
// That said this is a toy, and will remain a toy, so I can live with
// these problems.
//

package compiler

import (
	"fmt"
	"math"

	"github.com/skx/math-compiler/instructions"
	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

// Compiler holds our object-state.
type Compiler struct {

	// debug holds a flag to decide if debugging "stuff" is generated
	// in the output assembly.
	debug bool

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

//
// Our public API consists of the three functions:
//  New
//  SetDebug
//  Compile
//
// The rest of the code is an implementation detail.
//

// New creates a new compiler, given the expression in the constructor.
func New(input string) *Compiler {
	c := &Compiler{expression: input, constants: make(map[string]bool), debug: false}
	return c
}

// SetDebug changes the debug-flag for our output.
func (c *Compiler) SetDebug(val bool) {
	c.debug = val
}

// Compile converts the input program into a collection of
// AMD64-assembly language.
func (c *Compiler) Compile() (string, error) {

	//
	// Parse the program into a series of statements, etc.
	//
	// At this point there might be errors.  If so report them,
	// and terminate.
	//
	err := c.tokenize()
	if err != nil {
		return "", err
	}

	//
	// Convert the parsed-tokens to in internal-form.
	//
	c.makeinternalform()

	//
	// Now generate the output assembly
	//
	out := c.output()

	return out, nil
}

// tokenize populates our internal list of tokens, as a result of
// lexing the input string.
//
// There is some error-handling to ensure that the program looks
// somewhat reasonable.
func (c *Compiler) tokenize() error {

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
			return (fmt.Errorf("Error parsing input; token.ERROR returned from the lexer: %s", tok.Literal))
		}

		//
		// We'll convert "pi" and "e" into numbers as a special case.
		//
		if tok.Type == token.PI {
			tok.Type = token.NUMBER
			tok.Literal = fmt.Sprintf("%f", math.Pi)
		}
		if tok.Type == token.E {
			tok.Type = token.NUMBER
			tok.Literal = fmt.Sprintf("%f", math.E)
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
		return (fmt.Errorf("we expected the program to begin with a numeric thing"))
	}

	//
	// Get the last token
	//
	if len(c.tokens) > 1 {
		len := len(c.tokens)
		end := c.tokens[len-1]
		if end.Type == token.NUMBER {
			return fmt.Errorf("program ends with a number, which is invalid")
		}
	}

	//
	// No error.
	//
	return nil
}

// makeinternalform converts our series of tokens (i.e. the lexed expression)
// into an an intermediary form, collecting constants as they are discovered.
//
// This is the middle-step before generating our assembly-language program.
func (c *Compiler) makeinternalform() {

	//
	// Walk our tokens.
	//
	for _, t := range c.tokens {

		//
		// Handle each kind.
		//
		switch t.Type {

		case token.ABS:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Abs})

		case token.ASTERISK:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Multiply})

		case token.FACTORIAL:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Factorial})

		case token.COS:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Cos})

		case token.DUP:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Dup})

		case token.NUMBER:

			// Mark the constant as having been used.
			c.constants[t.Literal] = true

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Push, Value: t.Literal})

		case token.MOD:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Modulus})

		case token.MINUS:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Minus})

		case token.PLUS:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Plus})

		case token.POWER:

			// add the instruction
			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Power})

		case token.SIN:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Sin})

		case token.SLASH:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Divide})

		case token.SQRT:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Sqrt})

		case token.SWAP:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Swap})

		case token.TAN:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Tan})

		}
	}

}

// output generates the output, joining a header, a footer, and the
// writes our program to stdout
func (c *Compiler) output() string {

	//
	// The header.
	//
	header := `
#
# This assembly file was created by math-compiler.
#
# We're going to use Intel Syntax, because that is what I grew up with.
#
.intel_syntax noprefix
.global main

# Data-section:
#
# This contains data for the program at run-time.
#
#    int: is used to push values onto the stack.
#
#      a: used as an argument for functions that require one/two operands.
#
#      b: used as an argument for functions that require two operands.
#
#  depth: used to keep track of stack-depth.
#
#    fmt: Used to output the result of the calculation, later strings are for
#         various error-reports.
#
.data
          a: .double 0.0
          b: .double 0.0
      depth: .double 0.0
        int: .double 0.0

        fmt: .asciz "Result %g\n"
   div_zero: .asciz "Attempted division by zero.  Aborting\n"
   overflow: .asciz "Overflow - value out of range.  Aborting\n"
  stack_err: .asciz "Insufficient entries on the stack.  Aborting\n"
 stack_full: .asciz "Too many entries remaining on the stack.  Aborting\n"
`

	//
	// Output each of our discovered constants.
	//
	for v := range c.constants {
		header += fmt.Sprintf("%s: .double %s\n",
			c.escapeConstant(v), v)
	}

	header += `
#
# Main is our entry-point.
#
# We'll save rbp before we begin
#
main:
        push rbp

        # Our stack is initially empty (of numbers), so ensure that [depth]
        # is set to zero.
        #
        # Every time we push a value upon the stack we'll increase this value
        # and before we pop arguments from the stack we'll check there are
        # sufficient values stored.  This will prevent segfaults when user
        # programs are broken.
        #
        mov qword ptr [depth], 0

`
	if c.debug {
		header += "        # Debug-break\n"
		header += "        int 03\n"
	}

	//
	// The body of the program
	//
	body := ""

	// Now we walk over our internal-representation, and output
	// a chunk of assembly for each of our operator-types.
	for i, opr := range c.instructions {

		//
		// One-handler for each type: Alphabetical order.
		//
		switch opr.Type {

		case instructions.Abs:
			body += c.genAbs()

		case instructions.Cos:
			body += c.genCos()

		case instructions.Divide:
			body += c.genDivide()

		case instructions.Dup:
			body += c.genDup()

		case instructions.Factorial:
			body += c.genFactorial(i)

		case instructions.Minus:
			body += c.genMinus()

		case instructions.Modulus:
			body += c.genModulus()

		case instructions.Multiply:
			body += c.genMultiply()

		case instructions.Plus:
			body += c.genPlus()

		case instructions.Power:
			body += c.genPower(i)

		case instructions.Push:
			body += c.genPush(opr.Value)

		case instructions.Sin:
			body += c.genSin()

		case instructions.Sqrt:
			body += c.genSqrt()

		case instructions.Swap:
			body += c.genSwap()

		case instructions.Tan:
			body += c.genTan()

		}
	}

	footer := `
        # [PRINT]
        # ensure there is only one remaining argument upon the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jne stack_too_full      # should be only one entry.
        # print the result
        pop rax
        mov qword ptr [a], rax
        lea rdi,fmt             # format string
        movq xmm0, [a]          # argument
        movq rax, 1             # argument count
        call printf
        pop rbp
        xor rax,rax
        ret


#
# This is hit when a division by zero is attempted.
#
division_by_zero:
        lea rdi,div_zero
        jmp print_msg_and_exit

#
# This is hit when a register is too small to hold a value.
#
register_overflow:
        lea rdi,overflow
        jmp print_msg_and_exit


#
# This point is hit when the program is due to terminate, but the
# stack has too many entries upon it.
#
stack_too_full:
        lea rdi,stack_full
        jmp print_msg_and_exit

#
#
# This point is hit when there are insufficient operands upon the stack for
# a given operation.  (For example '3 +', or '3 4 + /'.)
#
stack_error:
        lea rdi,stack_err
        # jmp print_msg_and_exit - JMP is unnecessary here.

#
# Print a message and terminate.
#
# NOTE: We call 'exit' here to allow stdout to be flushed, and also to ensure
#       we don't need to balance our stack.
#
print_msg_and_exit:
        xor rax,rax
        call printf
        mov rdi,0
        call exit

`

	return header + body + footer
}
