package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

func main() {

	//
	// Setup the flags
	//
	compile := flag.Bool("compile", false, "Should we compile/execute by default?")
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
	var operations []string

	var i string

	for _, ent := range program {

		if i == "" {
			// number
			i = ent.Literal
		} else {

			// Number already set
			switch ent.Type {
			case token.PLUS:
				operations = append(operations, fmt.Sprintf("add rax, %s", i))
			case token.MINUS:
				operations = append(operations, fmt.Sprintf("sub rax,%s", i))
			case token.SLASH:
				// Look for the division by zero
				if i == "0" {
					operations = append(operations, "jmp div_by_zero")
				} else {
					operations = append(operations, fmt.Sprintf("mov rbx, %s", i))
					operations = append(operations, "cqo")
					operations = append(operations, fmt.Sprintf("div rbx"))
				}
			case token.ASTERISK:
				operations = append(operations, fmt.Sprintf("mov rbx, %s", i))
				operations = append(operations, fmt.Sprintf("mul rbx"))
			}
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
	out.Start = start
	out.Operations = operations

	//
	// This is the template we'll output.
	//
	assembly := `.intel_syntax noprefix
.global main

.data
format: .asciz "Division by zero\n"

main:
  mov rax, {{.Start}}
{{range .Operations}}  {{.}}
{{end}}
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
	// If we're not compiling then just output the assembly
	//
	if *compile == false {
		fmt.Printf(buf.String())
		os.Exit(0)
	}

	//
	// Get a sane value to write to
	//
	err = ioutil.WriteFile("tmp.s", buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing to tmp.s: %s\n", err.Error())
	}

	//
	// Now compile
	//
	_, err = exec.Command("gcc", "-static", "-o", "tmp.s.exe", "tmp.s").Output()

	//
	// Finally execute
	//
	cmd := exec.Command("tmp.s.exe")
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("%s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("%s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}
