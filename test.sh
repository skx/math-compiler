#!/bin/bash
#
# Simple test-driver to exercise our compiler.
#
# We compile some simple test-cases and test that the results match what
# we expect.
#


# Compile an expression and compare the result with a fixed value
#
# If the optional third argument is present, and non-empty, then we
# return the full output from the execution. Otherwise just the last
# token.
test_compile() {
    input="$1"
    result="$2"
    full="$3"

    #
    # Do this the long way round so we have assembly file for
    # inspection if/when a test fails.
    #
    rm -f test.s test || true
    go run main.go -- "${input}" > test.s
    gcc -static -o ./test test.s

    #
    # Run the test.
    #
    out=`./test`

    #
    # If we're not doing a "full" match we only take the last token
    #
    if [ "${full}" = "" ]; then
        out=$(echo "$out" | awk '{print $NF}')
    fi

    if [ "${result}" = "${out}" ]; then
        echo "Expected output found for '$input' [$result] "
        rm test test.s
    else
        echo "Expected output of '$input' is '$result' - got '${out}' instead"
        exit 1
    fi

}


# Simple operations
test_compile '1 2 3 4 + + +' 10
test_compile '3 4 5 * *'     60
test_compile '20 10 2 - -'   12
test_compile '20 4 2 / / '   10

# Division by zero
test_compile '3 4 /' '0.75'
test_compile '3 0 /' 'Attempted division by zero.  Aborting' 'full'
test_compile '3 5 5 - /' 'Attempted division by zero.  Aborting' 'full'

# Missing arguments upon the stack
test_compile '4 +' 'Insufficient entries on the stack.  Aborting' 'full'
test_compile '3 sin -' 'Insufficient entries on the stack.  Aborting' 'full'

# Too many arguments on the stack
test_compile '3 3 3 +' 'Too many entries remaining on the stack.  Aborting'  'full'


# modulus
test_compile  '1 4 %' 1
test_compile  '2 4 %' 2
test_compile  '3 4 %' 3
test_compile  '4 4 %' 0
test_compile  '5 4 %' 1
test_compile  '6 4 %' 2
test_compile  '7 4 %' 3
test_compile  '8 4 %' 0
test_compile  '9 4 %' 1
test_compile '10 4 %' 2
test_compile '11 4 %' 3
test_compile '12 4 %' 0

# powers of two - the manual-way
test_compile '2 2 *' 4
test_compile '2 2 2 * *' 8
test_compile '2 2 2 2 * * *' 16
test_compile '2 2 2 2 2 * * * *' 32
test_compile '2 2 2 2 2 2 * * * * *' 64
test_compile '2 2 2 2 2 2 2 * * * * * *' 128
test_compile '2 2 2 2 2 2 2 2 * * * * * * *' 256
test_compile '2 2 2 2 2 2 2 2 2 * * * * * * * *' 512
test_compile '2 2 2 2 2 2 2 2 2 2 * * * * * * * * *' 1024


# Add an extreme example of calculating 2^24:
inp="2 2 *"
for i in $(seq 1 22 ) ; do
    inp="${inp} 2 *"
done
test_compile "$inp" 1.67772e+07


# powers of two - the simple way
test_compile '2 0 ^'           0
test_compile '2 1 ^'           2
test_compile '2 2 ^'           4
test_compile '2 3 ^'           8
test_compile '2 4 ^'          16
test_compile '2 5 ^'          32
test_compile '2 6 ^'          64
test_compile '2 7 ^'         128
test_compile '2 8 ^'         256
test_compile '2 16 ^'      65536
test_compile '2 30 ^' 1.07374e+09


# division
test_compile '3 2 /' 1.5
test_compile '5 2 /' 2.5

# abs
test_compile '3 abs' 3
test_compile '3 9 - abs' 6
test_compile '-3 abs' 3

# sqrt
test_compile '9 sqrt' 3
test_compile '81 sqrt sqrt' 3
test_compile '81 sqrt sqrt sqrt' 1.73205


# circles
test_compile '1 sin' 0.841471
test_compile '1 cos' 0.540302
test_compile '1 tan' 1.55741

# swap
test_compile '3 5 -' -2
test_compile '3 5 swap -' 2

# dup
test_compile '3 sqrt dup *' 3
test_compile '3 dup ^' 27

exit 0
