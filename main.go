package main

import (
	"fmt"
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
	fmt.Printf("Starting with %s\n", start)

	//
	// Now process the rest of the program pair-wise
	//
	program = program[1:]

	for _, ent := range program {
		fmt.Printf("%v\n", ent)
	}
}
