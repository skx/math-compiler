// This is the main-driver for our compiler.

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
	debug := flag.Bool("debug", false, "Insert debug \"stuff\" in our generated output.")
	compile := flag.Bool("compile", false, "Compile the program, via invoking gcc.")
	program := flag.String("filename", "a.out", "The program to write to.")
	run := flag.Bool("run", false, "Run the binary, post-compile.")
	flag.Parse()

	//
	// If we're running we're also compiling
	//
	if *run == true {
		*compile = true
	}

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

	if *debug {
		comp.SetDebug(true)
	}

	//
	// Parse the program into a series of statements, etc.
	//
	// At this point there might be errors.  If so report them,
	// and terminate.
	//
	err := comp.Tokenize()
	if err != nil {
		fmt.Printf("There was an error compiling the input expression:\n")
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//
	// Convert the tokens to their internal form.
	//
	comp.InternalForm()

	//
	// Now generate the output assembly
	//
	var out string
	out, err = comp.Output()
	if err != nil {
		fmt.Printf("Error generating output from program:\n")
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	//
	// If we're not compiling then we just write the program to STDOUT.
	//
	if *compile == false {
		fmt.Printf("%s", out)
		return
	}

	//
	// OK we're compiling the program, via gcc.
	//
	gcc := exec.Command("gcc", "-static", "-o", *program, "-x", "assembler", "-")
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr

	//
	// We'll pipe our generated-program to STDIN of gcc, via a
	// temporary buffer-object.
	//
	var b bytes.Buffer
	b.Write([]byte(out))
	gcc.Stdin = &b

	//
	// Run gcc.
	//
	err = gcc.Run()
	if err != nil {
		fmt.Printf("Error launching gcc: %s\n", err)
		os.Exit(1)
	}

	//
	// Running the binary too?
	//
	if *run == true {
		exe := exec.Command(*program)
		exe.Stdout = os.Stdout
		exe.Stderr = os.Stderr
		err = exe.Run()
		if err != nil {
			fmt.Printf("Error launching %s: %s\n", *program, err)
			os.Exit(1)
		}
	}
}
