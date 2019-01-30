# math-compiler

This project contains the simplest possible compiler, which converts simple mathematical operations into assembly language, allowing all the speed in your sums!


# Installation

To install this glorious-tool, assuming you have a working golang installation:

    $ go get -u github.com/skx/math-compiler


## Quick Overview

The intention of this project is mostly to say "I wrote a compiler".

Because there are no shortages of toy-languages, and there is a lot of complexity in writing another for no real gain I decided to just focus upon the core:

* Allowing "maths" things to be "compiled".

In theory this would allow me to compile things like this:

    2 + ( 4 * 54 )

However I've even simplified that, via the use of reverse-polish notation.  So if you want to run that example you'd enter the expression as:

    4 54 * 2 +

## About Our Output

The output of this program will be an assembly-language file, which can be compiled and executed.

For example here is the simplest possible program:

    .intel_syntax noprefix
    .global main
    main:
        mov rax, 32
        ret

Given this program, saved in the file `test.s`, we can compile, then execute it like so:

     $ gcc -static -o test ./test.s
     $ ./test ; echo $?
     32



## Real Usage

Returning to our previous example of `2 + ( 4 * 54)` we can execute that via:

    $ math-compiler '4 54 * 2+' > sample.s
    $ gcc -static -o sample ./sample.s
    $ ./sample ; echo $?
    218

And you can compare that if you don't trust my maths (note that `*` is escaped to avoid your shell running a glob):

    $ expr 4 \* 54 + 2
    218

If you prefer to avoid compiling the output yourself you can do it directly, via the `-compile=true` argument:

    $ math-compiler -compile=true '3 34 + 3 *'
    111

## Test Cases

There are some test-cases contained in [test.sh](test.sh):

    $ ./test.sh
    Expected output found for '3 4 +' 7
    Expected output found for '3 4 *' 12
    Expected output found for '10 2 -' 8
    Expected output found for '10 2 /' 5


## Questions?

Great.  That concludes our exploration of compilers.

Steve
--
