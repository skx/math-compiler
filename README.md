# match-compiler

This project contains the simplest possible compiler, which converts
simple mathematical operations into assembly language, allowing
all the speed in your sums!

## Quick Overview

The intention of this project is mostly to say "I wrote a compiler".

Because there are no shortages of toy-languages, and there is a lot of
complexity in writing another for no real gain I decided to just focus
upon the core:

* Allowing "maths" things to be "compiled".

In theory this would allow me to compile things like this:

    2 + 4 * 54

However I've even simplified that, via the use of reverse-polish notation.


## Assembly Output

The output of this program will be an assembly file.  For example here is
the simplest possible program:


        .intel_syntax noprefix
        .global main
    main:
        mov rax, 32
        ret

Given this program, in the file `test.s` we can compile it, and execute it like so:

     $ gcc -static -o test ./test.s
     $ ./test ; echo $?
     32


## Examples

TODO
