// compiler contains our simple compiler.
//
// In brief it uses the lexer to tokenize the expression, then we convert
// that series of tokens into an "internal representation" which is pretty
// much things like:
//
//    push_int 3
//    push_int 4
//    Add
//
// We iterate over this simple representation and output a block of code
// for each.
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
	"strconv"
	"strings"

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

// New creates a new compiler, given the expression in the constructor.
func New(input string) *Compiler {
	c := &Compiler{expression: input, constants: make(map[string]bool), debug: false}
	return c
}

// SetDebug changes the debug-flag for our output.
func (c *Compiler) SetDebug(val bool) {
	c.debug = val
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
// an an intermediary form, collecting constants as they are discovered.
//
// This is the middle-step before generating our assembly-language program.
func (c *Compiler) InternalForm() {

	//
	// Walk our tokens.
	//
	for _, t := range c.tokens {

		//
		// Handle each kind.
		//
		switch t.Type {

		case token.ASTERISK:

			c.instructions = append(c.instructions,
				instructions.Instruction{Type: instructions.Multiply})

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

// Output writes our program to stdout
func (c *Compiler) Output() (string, error) {

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
#    fmt: Used to output the result of the calculation, and div_zero for
#         the obvious error-case.
#
.data
          a: .double 0.0
          b: .double 0.0
      depth: .double 0.0
        int: .double 0.0
        fmt: .asciz "Result %g\n"
   div_zero: .asciz "Attempted division by zero.  Aborting\n"
`

	//
	// Add on the constants
	//
	for v, _ := range c.constants {
		header += fmt.Sprintf("%s: .double %s\n",
			c.escapeConstant(v), v)
	}

	header += `
#
# Main is our entry-point.
#
# We'll save rbp before we begin
main:
        push rbp

        # Our stack is initially empty (of numbers), so ensure that [depth]
        # is set to zero.
        # Every time we push a value upon the stack we'll increase this value
        # and before we pop arguments from the stack we'll check there are
        # sufficient values stored.  This will prevent segfaults when user
        # programs are broken.
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

		case instructions.Cos:
			body += c.genCos()

		case instructions.Divide:
			body += c.genDivide()

		case instructions.Dup:
			body += c.genDup()

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
        # Now we pop the last item from the stack, and print the result
        # then we terminate cleanly.
        pop rax
        mov qword ptr [a], rax
        movq xmm0, [a]
        lea rdi,fmt
        movq rax, 1
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
# This point is hit when there are insufficient operands upon the stack for
# a given operation.  (For example '3 +', or '3 4 + /'.)
#
stack_error:
        lea rdi,stack_err
        # jmp print_msg_and_exit - JMP is unnecessary here.

#
# Print a message and terminate.
#
# NOTE: We call 'exit' here to allow stdout to be flushed.
#
print_msg_and_exit:
        xor rax, rax
        call printf
        mov rdi,0
        call exit

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

	// remove periods
	val = strings.Replace(val, ".", "_", -1)
	// remove minus-signs
	val = strings.Replace(val, "-", "", -1)
	return val
}

// genCos generates assembly code to pop a value from the stack,
// run a cos-operation, and store the result back on the stack.
func (c *Compiler) genCos() string {
	return `
        # pop one value
        pop rax
        mov qword ptr [a], rax

        # cos
        fld qword ptr [a]
        fcos
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`
}

// genDivide generates assembly code to pop two values from the stack,
// divide them and store the result back on the stack.
func (c *Compiler) genDivide() string {
	return `
        # pop two values
        pop rax
        cmp rax,0
        je division_by_zero
        mov qword ptr [a], rax
        pop rax
        mov qword ptr [b], rax

        # divide
        fld qword ptr [b]
        fdiv qword ptr  [a]
        fstp qword ptr [a]

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax
`

}

// genDup generates assembly code to pop a value from the stack and
// push it back twice - effectively duplicating it.
func (c *Compiler) genDup() string {
	return `
        pop rax
        push rax
        push rax
`
}

// genMinus generates assembly code to pop two values from the stack,
// subtract them and store the result back on the stack.
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

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax
`
}

// genModulus generates assembly code to pop two values from the stack,
// perform a modulus-operation and store the result back on the stack.
// Note we truncate things to integers in this section of the code.
func (c *Compiler) genModulus() string {
	return `

        # pop two values - rounding both to ints
        pop rax
        mov qword ptr [a], rax
        fld qword ptr [a]
        frndint
        fistp qword ptr [a]

        pop rax
        mov qword ptr [b], rax
        fld qword ptr [b]
        frndint
        fistp qword ptr [b]

        # now we do the modulus-magic.
        mov rax, qword ptr [b]
        mov rbx, qword ptr [a]
        xor rdx, rdx
        cqo
        div rbx

        # store the result from 'rdx'.
        mov qword ptr[a], rdx
        fild qword ptr [a]
        fstp qword ptr [a]
        mov rax, qword ptr [a]
        push rax
`
}

// genMultiply generates assembly code to pop two values from the stack,
// multiply them and store the result back on the stack.
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

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax
`

}

// genPlus generates assembly code to pop two values from the stack,
// add them and store the result back on the stack.
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

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax
`
}

// genPower generates assembly code to pop two values from the stack,
// perform a power-raising and store the result back on the stack.
//
// Note we truncate things to integers in this section of the code.
//
// Note we do some comparisions here, and need to generate some (unique) labels
//
func (c *Compiler) genPower(i int) string {
	text := `

        # pop two values - rounding both to ints
        pop rax
        mov qword ptr [a], rax
        fld qword ptr [a]
        frndint
        fistp qword ptr [a]

        pop rax
        mov qword ptr [b], rax
        fld qword ptr [b]
        frndint
        fistp qword ptr [b]

        # get the two values
        mov rax, qword ptr [b]
        mov rbx, qword ptr [a]

        # if the power is 0 we return zero
        cmp rbx, 0
        jne none_zero_#ID
           # store zero
           fldz
           jmp store_value_#ID

none_zero_#ID:

        # if the power is 1 we return the original value
        cmp rbx, 1
        jne none_one_#ID
           mov qword ptr[a], rax
           fild qword ptr [a]
           jmp store_value_#ID

none_one_#ID:
        # here we have rax having a value
        # and we have rbx having the power to raise
        mov rcx, rax   # save the value

        # decrease the power by one.
        dec rbx
again_#ID:
           # rax = rax * rcx (which is the original value we started with)
           imul rax,rcx
           dec rbx
           jnz again_#ID

        mov qword ptr[a], rax
        fild qword ptr [a]

store_value_#ID:

        fstp qword ptr [a]

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax
`

	return (strings.Replace(text, "#ID", fmt.Sprintf("%d", i), -1))
}

// genPush generates assembly code to push a value upon the RPN stack.
func (c *Compiler) genPush(value string) string {

	text := `
        fld qword ptr #VAL
        fstp qword ptr [int]
        mov rax, qword ptr [int]
        push rax
`

	return (strings.Replace(text, "#VAL", c.escapeConstant(value), -1))
}

// genSin generates assembly code to pop a value from the stack,
// run a sin-operation, and store the result back on the stack.
func (c *Compiler) genSin() string {
	return `
        # pop one value
        pop rax
        mov qword ptr [a], rax

        # sin
        fld qword ptr [a]
        fsin
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`
}

// genSwap generates assembly code to pop two values from the stack and
// push them back, in the other order.
func (c *Compiler) genSwap() string {
	return `
        pop rax
        pop rbx
        push rax
        push rbx
`
}

// genTan generates assembly code to pop a value from the stack,
// run a tan-operation, and store the result back on the stack.
func (c *Compiler) genTan() string {
	return `
        # pop one value
        pop rax
        mov qword ptr [a], rax

        # tan
        fld qword ptr [a]
        fsincos
        fdivr %st(0),st(1)
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`
}

// genSqrt generates assembly code to pop a value from the stack,
// run a square-root operation, and store the result back on the stack.
func (c *Compiler) genSqrt() string {
	return `
        # pop one value
        pop rax
        mov qword ptr [a], rax

        # sqrt
        fld qword ptr [a]
        fsqrt
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax
`
}
