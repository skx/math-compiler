package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

func main() {

	//
	// Look for flags.
	//
	compile := flag.Bool("compile", false, "Compile the program, to a.out")
	flag.Parse()

	//
	// Ensure we have a single argument
	//
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: math-compiler 'expression'\n")
		os.Exit(1)
	}

	//
	// Create the lexer - based upon our argument
	//
	input := flag.Args()[0]
	lexed := lexer.New(input)

	//
	// Now we want to process the tokens.
	//
	// Because we're using reverse-polish we'll assume
	// we have input of the form:
	//
	//  3 4 + 5 *
	//
	// (This represents the input '(3 + 4) * 5)').
	//
	// So we have an initial digit, and we have then
	// pairs of ["number" "operation"]
	//
	var program []token.Token

	//
	// First of all populate that `program` array with our tokens.
	//
	for {
		tok := lexed.NextToken()
		if tok.Type == token.EOF {
			break
		}
		program = append(program, tok)
	}

	//
	// If the first token isn't a number we're in trouble
	//
	if len(program) < 1 {
		fmt.Printf("Empty program!\n")
		os.Exit(1)
	}
	if program[0].Type != token.INT {
		fmt.Printf("We expected the program to begin with an integer!\n")
		os.Exit(1)
	}

	//
	// Pop off the starting integer.
	//
	start := program[0].Literal

	//
	// Now process the rest of the program pair-wise - removing that
	// first number.
	//
	program = program[1:]

	//
	// We expect to work in pairs; reading two elements from
	// our program.
	//
	// For example with the program:
	//
	//  3 4 + 5 * 2 /
	//
	// We've already removed the leading "3" so now we expect
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
	// TODO: We could abort here if the length of `program` was not even.
	//
	var operations []string

	var i string

	for offset, ent := range program {

		if i == "" {
			// number
			i = ent.Literal
		} else {

			//
			// The number is already set, so we're now expecting
			// an operator.
			//
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
					// means we need a uniq label
					//
					// We generate the label ID as the offset of
					// the statement we're generating in our input
					//
					operations = append(operations, `mov rcx, `+i)
					operations = append(operations, `mov ebx, eax`)
					operations = append(operations, `dec rcx`)
					operations = append(operations, fmt.Sprintf("label_%d:", offset))
					operations = append(operations, `  mul ebx`)
					operations = append(operations, `  dec rcx`)
					operations = append(operations, fmt.Sprintf("  jnz label_%d", offset))
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
				fmt.Printf("Invalid program - expected operator, but found %v\n", ent)
				os.Exit(1)
			}

			//
			// Next time around the loop we'll be looking for
			// a number, rather than an operator.
			//
			i = ""
		}
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
		fmt.Printf("Error compiling template: %s\n", err.Error())
		os.Exit(1)
	}

	//
	// Finally show that to STDOUT
	//
	if *compile == false {

		fmt.Printf("%s", buf.String())
	} else {
		gcc := exec.Command("gcc", "-static", "-o", "a.out", "-x", "assembler", "-")
		gcc.Stdout = os.Stdout
		gcc.Stdin = buf
		gcc.Stderr = os.Stderr

		err := gcc.Run()
		if err != nil {
			fmt.Printf("Error launching gcc: %s\n", err)
			os.Exit(1)
		}

	}
}
