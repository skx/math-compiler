# math-compiler

This project contains the simplest possible compiler, which converts simple mathematical operations into assembly language, allowing all the speed in your sums!



# Installation

To install this glorious-tool, assuming you have a working [golang](https://golang.org/) installation:

    $ go get -u github.com/skx/math-compiler



## Quick Overview

The intention of this project is mostly to say "I wrote a compiler", because I've already [experimented with a language](https://github.com/skx/monkey/), and [implemented a BASIC](https://github.com/skx/gobasic/).  The things learned from those projects were pretty useful, even if the actual results were not so obviously useful in themselves.

Because there are no shortages of toy-languages, and there is a lot of complexity in writing another for no real gain, I decided to just focus upon a simple core:

* Allowing "maths stuff" to be "compiled".

In theory this would allow me to compile things like this:

    2 + ( 4 * 54 )

However I even simplified that, via the use of a "[Reverse Polish](https://en.wikipedia.org/wiki/Reverse_Polish_notation)"-like notation, so if you want to run that example you'd enter the expression as:

    4 54 * 2 +

(In our case we have no stack, so we just assume "operand ... operand operator", constantly.)


## Supported Operations

We support the following operations:

* `+` - Plus
  * `1 + 2` is expressed `1 2 +`
* `-` - Minus
  * `4 - 2` is expressed `4 2 -`
* `*` - Multiply
  * `4 * 4` is expressed `4 4 *`
* `/` - Divide
  * `40 / 10` is expressed `40 10 /`
* `^` - Raise to a power
  * `2 ^ 4` is expressed `2 4 ^`
* `%` - Modulus
  * `7 % 5` is expressed `7 5 %`

Hopefully you can see how those are chained, so the sum `(3 + (5 * 5)) / 2` would be written as `5 5 * 3 + 2 /`.



## About Our Output

The output of `math-compiler` will typically be an assembly-language file, which then needs to be compiled before it may be executed.

Returning to our previous example of `2 + ( 4 * 54)` we can compile & execute that program like so:

    $ math-compiler '4 54 * 2+' > sample.s
    $ gcc -static -o sample ./sample.s
    $ ./sample
    Result 218

There you see:

* `math-compiler` was invoked, and the output written to the file `sample.s`.
* `gcc` was used to assemble `sample.s` into the binary `sample`.
* The actual binary was then executed, and that showed the result.

If you prefer you can also let the compiler do the heavy-lifting, and generate an executable for you directly.  Simply add `-compile`, and execute the generated `a.out` binary:

    $ math-compiler -compile=true '2 8 ^'
    $ a.out
    Result 256



## Test Cases

The codebase itself contains some simple test-cases, however these are
not comprehensive.

To execute the tests use the standard approach:

    $ go test [-race] ./...

In addition to the internal test cases there are also some functional tests
contained in [test.sh](test.sh) - these perform some calculations and verify
they produce the correct result.

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

As you can see the registers can store a different number of bits, depending on how much you access:

     0x1122334455667788
     ================ rax (64 bits)
             ======== eax (32 bits)
                 ====  ax (16 bits)
                   ==  ah (8 bits)
                   ==  al (8 bits)

I believe that means we should be OK to store 64-bit numbers.



## Questions?

Great.  That concludes our exploration of compilers.



Steve
--
