// generator.go contains the code for emitting instructions.

package compiler

import (
	"fmt"
	"strconv"
	"strings"
)

// escapeConstant converts a floating-point number such as
// "1.2", or "-1.3" into a constant value that can be embedded
// safely into our generated assembly-language file.
func (c *Compiler) escapeConstant(input string) string {

	// Convert "3.0" to "const_3.0"
	s, _ := strconv.ParseFloat(input, 32)

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

// genAbs generates assembly code to pop a value from the stack,
// run an ABS-operation, and store the result back on the stack.
func (c *Compiler) genAbs() string {
	return `
        # [ABS]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

        # pop one value
        pop rax
        mov qword ptr [a], rax

        # abs
        fld qword ptr [a]
        fabs
        fstp qword ptr [a]

        # push result onto stack
        mov rax, qword ptr [a]
        push rax

        # stack size didn't change; popped one, pushed one.
`
}

// genCos generates assembly code to pop a value from the stack,
// run a cos-operation, and store the result back on the stack.
func (c *Compiler) genCos() string {
	return `
        # [COS]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

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

        # stack size didn't change; popped one, pushed one.
`
}

// genDivide generates assembly code to pop two values from the stack,
// divide them and store the result back on the stack.
func (c *Compiler) genDivide() string {
	return `
        # [DIVIDE]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`

}

// genDup generates assembly code to pop a value from the stack and
// push it back twice - effectively duplicating it.
func (c *Compiler) genDup() string {
	return `
        # [DUP]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

        pop rax
        push rax
        push rax

        # We've added a new entry to the stack.
        inc qword ptr [depth]
`
}

// genFactorial generates assembly code to pop a value from the stack,
// run a factorial-operation, and store the result back on the stack.
func (c *Compiler) genFactorial(i int) string {
	text := `
        # [FACTORIAL]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

        # pop a value - rounding to an int
        pop rax
        mov qword ptr [a], rax
        fld qword ptr [a]
        frndint
        fistp qword ptr [a]

        # get the value in rcx, setup rax to be 1
        mov rcx, qword ptr [a]
        mov rax,1

        # If the value is negative, return zero
        cmp rcx, 0
        jg again_#ID
        # jg means jump-if-greater, so if we hit this we had zero/negative
        # store the result.
        mov qword ptr[a], 0
        jmp store_result_#ID

again_#ID:
        # rax = rax * rcx
        imul rax, rcx

        # value too big?
        jo register_overflow

        dec rcx
        jnz again_#ID

        # store
        mov qword ptr[a], rax
        fild qword ptr [a]
        fstp qword ptr [a]
        mov rax, qword ptr [a]

store_result_#ID:
        # push result onto stack
        mov rax, qword ptr [a]
        push rax
        # stack size didn't change; popped one, pushed one.
`
	return (strings.Replace(text, "#ID", fmt.Sprintf("%d", i), -1))
}

// genMinus generates assembly code to pop two values from the stack,
// subtract them and store the result back on the stack.
func (c *Compiler) genMinus() string {
	return `
        # [MINUS]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`
}

// genModulus generates assembly code to pop two values from the stack,
// perform a modulus-operation and store the result back on the stack.
// Note we truncate things to integers in this section of the code.
func (c *Compiler) genModulus() string {
	return `
        # [MODULUS]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`
}

// genMultiply generates assembly code to pop two values from the stack,
// multiply them and store the result back on the stack.
func (c *Compiler) genMultiply() string {
	return `
        # [MULTIPLY]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`

}

// genPlus generates assembly code to pop two values from the stack,
// add them and store the result back on the stack.
func (c *Compiler) genPlus() string {

	return `
        # [PLUS]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`
}

// genPower generates assembly code to pop two values from the stack,
// perform a power-raising and store the result back on the stack.
//
// Note we truncate things to integers in this section of the code.
//
// Note we do some comparisons here, and need to generate some (unique) labels
//
func (c *Compiler) genPower(i int) string {
	text := `
        # [POWER]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

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

           # value too big?
           jo register_overflow

           dec rbx
           jnz again_#ID

        mov qword ptr[a], rax
        fild qword ptr [a]

store_value_#ID:

        fstp qword ptr [a]

        # push the result back onto the stack
        mov rax, qword ptr [a]
        push rax

        # we took two values from the stack, but added one
        # so the net result is the stack shrunk by one.
        dec qword ptr [depth]
`

	return (strings.Replace(text, "#ID", fmt.Sprintf("%d", i), -1))
}

// genPush generates assembly code to push a value upon the RPN stack.
func (c *Compiler) genPush(value string) string {

	text := `
        # [PUSH]
        # Load the value #VALUE onto the stack
        # Increase the value stored at [depth] to note we've a new stack-entry
        fld qword ptr #ESCAPED
        fstp qword ptr [int]
        mov rax, qword ptr [int]
        push rax
        inc qword ptr [depth]
`

	// Allow the value and the escaped value to be expanded.
	text = strings.Replace(text, "#VALUE", value, -1)
	text = strings.Replace(text, "#ESCAPED", c.escapeConstant(value), -1)

	return (text)
}

// genSin generates assembly code to pop a value from the stack,
// run a sin-operation, and store the result back on the stack.
func (c *Compiler) genSin() string {
	return `
        # [SIN]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

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

        # stack size didn't change; popped one, pushed one.
`
}

// genSwap generates assembly code to pop two values from the stack and
// push them back, in the other order.
func (c *Compiler) genSwap() string {
	return `
        # [SWAP]
        # ensure there are at least two arguments on the stack
        mov rax, qword ptr [depth]
        cmp rax, 2
        jb stack_error

        pop rax
        pop rbx
        push rax
        push rbx
        # stack size didn't change; popped two, pushed two.
`
}

// genTan generates assembly code to pop a value from the stack,
// run a tan-operation, and store the result back on the stack.
func (c *Compiler) genTan() string {
	return `
        # [TAN]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

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
        # [SQRT]
        # ensure there is at least one argument on the stack
        mov rax, qword ptr [depth]
        cmp rax, 1
        jb stack_error

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

        # stack size didn't change; popped one, pushed one.
`
}
