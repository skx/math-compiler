[![Travis CI](https://img.shields.io/travis/skx/math-compiler/master.svg?style=flat-square)](https://travis-ci.org/skx/puppet-summary)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/math-compiler)](https://goreportcard.com/report/github.com/skx/puppet-summary)
[![license](https://img.shields.io/github/license/skx/math-compiler.svg)](https://github.com/skx/puppet-summary/blob/master/LICENSE)

# math-compiler

This project contains the simplest possible compiler, which converts mathematical operations into assembly language, allowing all the speed in your sums!

Because this is a simple project it provides only a small number of primitives:

* `+` - Plus
* `-` - Minus
* `*` - Multiply
* `/` - Divide
* `^` - Raise to a power
* `%` - Modulus
* `abs`
* `sin`
* `cos`
* `tan`
* `sqrt`
* Stack operations:
  * `swap` - Swap the top-two items on the stack
  * `dup` - Duplicate the topmost stack-entry.
* Built-in constants:
  * `e`
  * `pi`

Despite this toy-functionality there is a lot going on, and we support:

* Full RPN input
* Floating-point numbers (i.e. one-third multipled by nine is 3)
   * `1 3 / 9 *`
* Negative numbers work as you'd expect.

Some errors will be caught at run-time, as the generated code has support for:

* Detecting, and preventing, division by zero.
* Detecting insufficient arguments being present upon the stack.
  * For example this program is invalid `3 +`, because the addition operator requires two operands.  (i.e. `3 4 +`)


# Installation


Providing you have a working [go-installation](https://golang.org/) you should be able to install this software by running:

    $ go get -u github.com/skx/math-compiler



## Quick Overview

The intention of this project is mostly to say "I wrote a compiler", because I've already [experimented with a language](https://github.com/skx/monkey/), and [implemented a BASIC](https://github.com/skx/gobasic/).  The things learned from those projects were pretty useful, even if the actual results were not so obviously useful in themselves.

Because there are no shortages of toy-languages, and there is a lot of complexity in writing another for no real gain, I decided to just focus upon a simple core:

* Allowing "maths stuff" to be "compiled".

In theory this would allow me to compile things like this:

    2 + ( 4 * 54 )

However I even simplified that, via the use of a "[Reverse Polish](https://en.wikipedia.org/wiki/Reverse_Polish_notation)" notation, so if you want to run that example you'd enter the expression as:

    4 54 * 2 +




## About Our Output

The output of `math-compiler` will typically be an assembly-language file, which then needs to be compiled before it may be executed.

Given our previous example of `2 + ( 4 * 54)` we can compile & execute that program like so:

    $ math-compiler '4 54 * 2+' > sample.s
    $ gcc -static -o sample ./sample.s
    $ ./sample
    Result 218

There you see:

* `math-compiler` was invoked, and the output written to the file `sample.s`.
* `gcc` was used to assemble `sample.s` into the binary `sample`.
* The actual binary was then executed, which showed the result of the calculation.

If you prefer you can also let the compiler do the heavy-lifting, and generate an executable for you directly.  Simply add `-compile`, and execute the generated `a.out` binary:

    $ math-compiler -compile=true '2 8 ^'
    $ ./a.out
    Result 256

Or to compile __and__ execute directly:

    $ math-compiler -run '3 45 * 9 + 12 /'
    Result 12


## Test Cases

The codebase itself contains some simple test-cases, however these are not comprehensive as a large part of our operation is merely to populate a simple template-file, and it is hard to test that.

To execute the integrated tests use the standard go approach:

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




### Debugging the generated programs

If you run the compiler with the `-debug` flag a breakpoint will be generated
immediately at the start of the program.  You can use that breakpoint to easily
debug the generated binary via `gdb`.

For example you might generate a program "`2 3 + 4 /`" like so:

    $ math-compiler -compile -debug '2 3 + 4 /'

Now you can launch that binary under `gdb`, and run it:

    $ gdb ./a.out
    (gdb) run
    ..
    Program received signal SIGTRAP, Trace/breakpoint trap.
    0x00000000006b20cd in main ()

Dissassemble the code via `disassemble`, and step over instructions one at a time via `stepi`.  If your program is long you might see a lot of output from the `disassemble` step:

    (gdb) disassemble
    Dump of assembler code for function main:
       0x00000000006b20cb:	push   %rbp
       0x00000000006b20cc:	int3
    => 0x00000000006b20cd:	fldl   0x6b20b3
       0x00000000006b20d4:	fstpl  0x6b2090
       0x00000000006b20db:	mov    0x6b2090,%rax
       0x00000000006b20e3:	push   %rax
       0x00000000006b20e4:	fldl   0x6b20bb
       0x00000000006b20eb:	fstpl  0x6b2090
       0x00000000006b20f2:	mov    0x6b2090,%rax
       0x00000000006b20fa:	push   %rax
       ...
       ...

You can set a breakpoint at a line in the future, and continue running till
you hit it, with something like this:

     (gdb) break *0x00000000006b20fa
     (gdb) cont

Once there inspect the registers with commands like these two:

     (gdb) print $rax
     (gdb) info registers

My favourite is `info registers float`, which shows you the floating-point
values as well as the raw values:

     (gdb) info registers float
     st0            0.140652076786443369638	(raw 0x3ffc90071917a6263000)
     st1            0	(raw 0x00000000000000000000)
     st2            0	(raw 0x00000000000000000000)
     ...
     ...

Further documentation can be found in the `gdb` manual, which is worth reading
if you've an interest in compilers, debuggers, and decompilers.


## Possible Expansion?

The obvious thing to improve in this compiler is to add support for more floating-point operations.  At the moment basic-support is present, allowing calcuations such as this to produce the correct result:

* `3 2 /`
  * Correctly returns `1.5`
* `1 3 / 9 *`
  * Correctly returns 1/3 * 9 == `3`.
* `81 sqrt sqrt`
  * Correctly returns `root(root(81))`



## Questions?

Great.  That concludes our exploration of compilers.



Steve
--
