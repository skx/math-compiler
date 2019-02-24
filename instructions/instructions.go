// Package instructions contains a series of types.
//
// We parse the program we're given into a series of tokens, then
// later convert those tokens into an internal-form, using these
// instructions.
//
// At generation time we emit a header, a footer, and a series of
// snippets - one for each logical instruction.
package instructions

// InstructionType holds the type of the instruction.
type InstructionType byte

const (
	// Push is used to generate code to push a number onto the stack.
	Push InstructionType = 'p'

	// Plus means to pop two items from the stack and push the result
	// of adding them.
	Plus InstructionType = '+'

	// Minus means to pop two items from the stack and push the result
	// of subtracting them.
	Minus InstructionType = '-'

	// Multiply means to pop two items from the stack and push the result
	// of multiplying them.
	Multiply InstructionType = '*'

	// Divide means to pop two items from the stack and push the result
	// of dividing them.
	Divide InstructionType = '/'

	// Power means to pop two items from the stack and push the result
	// of raising one to the power of the other.
	Power InstructionType = '^'

	// Modulus means to pop two items from the stack and push the result
	// of running a modulus operation.
	Modulus InstructionType = '%'

	// Abs is used to pop a value from the stack and push the absolute
	// value back.
	Abs InstructionType = 'a'

	// Sin is used to pop a value from the stack and push the result
	// of sin() back.
	Sin InstructionType = 's'

	// Cos is used to pop a value from the stack and push the result
	// of cos() back.
	Cos InstructionType = 'c'

	// Tan is used to pop a value from the stack and push the result
	// of tan() back.
	Tan InstructionType = 't'

	// Sqrt is used to pop a value from the stack and push the result
	// of calculating its square-root back.
	Sqrt InstructionType = 'q'

	// Swap swaps the position of the top two stack-items.
	Swap InstructionType = 'S'

	// Dup duplicates the stacks topmost value.
	Dup InstructionType = 'D'
)

// Instruction holds a single thing that the compiler must generate code for.
// (The value is only used when a float is to be pushed upon the stack.)
type Instruction struct {

	// Type holds the type of instruction this object represents
	Type InstructionType

	// Value holds the value of a number to be pushed upon the RPN stack.
	Value string
}
