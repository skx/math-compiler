package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/skx/math-compiler/lexer"
	"github.com/skx/math-compiler/token"
)

func main() {

	//
	// Ensure we have a single argument
	//
	if len(os.Args) != 2 {
		fmt.Printf("Usage: math-compiler 'expression'\n")
		os.Exit(1)
	}

	//
	// Create the lexer - based upon our argument
	//
	input := os.Args[1]
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
				operations = append(operations, fmt.Sprintf("mov rbx, %s", i))
				operations = append(operations, "cqo")
				operations = append(operations, fmt.Sprintf("div rbx"))
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

main:
  mov rax, {{.Start}}
{{range .Operations}}  {{.}}
{{end}}
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

	//
	// If there are no errors then write it to the console
	//
	if err == nil {
		fmt.Printf(buf.String())
	} else {
		fmt.Printf("Error compiling template: %s\n", err.Error())
	}

}
