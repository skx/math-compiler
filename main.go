// This is the main-driver for our compiler.
//

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/skx/math-compiler/compiler"
)

func main() {

	//
	// Look for flags.
	//
	compile := flag.Bool("compile", false, "Compile the program, to a.out")
	flag.Parse()

	//
	// Ensure we have an expression as our single argument.
	//
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: math-compiler 'expression'\n")
		os.Exit(1)
	}

	//
	// Create a compiler-object, with the program as input.
	//
	comp := compiler.New(flag.Args()[0])

	//
	// Parse the program into a series of statements, etc.
	//
	// At this point there might be errors.  If so report that.
	//
	err := comp.Compile()
	if err != nil {
		fmt.Printf("There was an error compiling the input expression:\n")
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//
	// Now generate a program from the expression.
	//
	var out string
	out, err = comp.Output()
	if err != nil {
		fmt.Printf("Error generating output from program:\n")
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//
	// If we're not compiling then we just write that to STDOUT.
	//
	if *compile == false {
		fmt.Printf("%s", out)
		return
	}

	//
	// OK we're compiling the program directly into the file
	// `a.out`.
	//
	// Do that by invoking gcc.
	//
	bin := "a.out"
	gcc := exec.Command("gcc", "-static", "-o", bin, "-x", "assembler", "-")
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr

	//
	// We'll pipe our generated-program to STDIN of gcc, via an
	// interim-buffer object.
	//
	var b bytes.Buffer
	b.Write([]byte(out))
	gcc.Stdin = &b

	err = gcc.Run()
	if err != nil {
		fmt.Printf("Error launching gcc: %s\n", err)
		os.Exit(1)
	}
}
