// instruction.go
package instructions

// The type of the instruction
type InstructionType byte

const (
	// setup: push number on stack
	Push InstructionType = 'p'

	// simple operators that work with stack-numbers
	Plus     InstructionType = '+'
	Minus    InstructionType = '-'
	Multiply InstructionType = '*'
	Divide   InstructionType = '/'

	// medium operators that work with stack-numbers
	Power   InstructionType = '^'
	Modulus InstructionType = '%'

	// complex operators that work with stack-numbers
	Sin  InstructionType = 's'
	Cos  InstructionType = 'c'
	Tan  InstructionType = 't'
	Sqrt InstructionType = 'q'
)

// A single instruction will have a thing to do, and optionally a value.
//
// (The value is only used when a float is to be pushed upon the stack.)
type Instruction struct {

	// The type of instruction this object represents
	Type InstructionType

	// Value holds the  value of a number to be pushed upon the RPN stack.
	Value string
}
