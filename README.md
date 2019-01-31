# math-compiler

This project contains the simplest possible compiler, which converts simple mathematical operations into assembly language, allowing all the speed in your sums!


# Installation

To install this glorious-tool, assuming you have a working [golang](https://golang.org/) installation:

    $ go get -u github.com/skx/math-compiler


## Quick Overview

The intention of this project is mostly to say "I wrote a compiler", because I've already [experimented with a language](https://github.com/skx/monkey/), and [implemented a BASIC](https://github.com/skx/gobasic/).  The things learned from this were pretty useful, even if the end-results were not so obviously useful.

Because there are no shortages of toy-languages, and there is a lot of complexity in writing another for no real gain I decided to just focus upon the core:

* Allowing "maths" things to be "compiled".

In theory this would allow me to compile things like this:

    2 + ( 4 * 54 )

However I've even simplified that, via the use of a "[Reverse Polish](https://en.wikipedia.org/wiki/Reverse_Polish_notation)"-like notation.  So if you want to run that example you'd enter the expression as:

    4 54 * 2 +

(In our case we have no stack, so we just assume "operand ... operand operator", constantly.)


## Supported Operations

We support the following operations:

* `+` - Plus
* `-` - Minus
* `*` - Multiply
* `/` - Divide
* `^` - Raise to a power
* `%` - Modulus



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
    $ ./sample
    Result 218

And you can compare that if you don't trust my maths (note that `*` is escaped to avoid your shell running a glob):

    $ expr 4 \* 54 + 2
    218

If you wish you can "simplify" the compilation step by piping the output directly to `gcc`, you'll need to add `-x assembler` because the filename-suffix will no longer be available to allow language-detection.

Here is an example:

    $ ./math-compiler '6 6 * 12 /' | gcc -static -o prog  -x assembler - ; ./prog
    Result 3

Or failing that you can allow the `math-compiler` to generate an executable for you:

    $ ./math-compiler -compile=true '2 8 ^'
    $ ./a.out
    Result 256

## Test Cases

There are some test-cases contained in [test.sh](test.sh):

    frodo ~/go/src/github.com/skx/math-compiler $ ./test.sh
    ...
    Expected output found for '2 0 ^' [0]
    Expected output found for '2 1 ^' [2]
    Expected output found for '2 2 ^' [4]
    Expected output found for '2 3 ^' [8]
    Expected output found for '2 4 ^' [16]
    Expected output found for '2 5 ^' [32]
    ...
    Expected output found for '2 30 ^' [1073741824]
    ...



## Numerical Limits

I try to use full-width instructions where possible.

As ou can see the intel registers can store a different number of bits:

     0x1122334455667788
     ================ rax (64 bits)
             ======== eax (32 bits)
                 ====  ax (16 bits)
                   ==    ah (8 bits)
                   ==  al (8 bits)

I believe that means we should be OK to store 64-bit numbers.


## Questions?

Great.  That concludes our exploration of compilers.

Steve
--
