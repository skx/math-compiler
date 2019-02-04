#!/bin/sh
#
# Simple test-driver to exercise our compiler.
#
# We compile some simple test-cases and test that the results match what
# we expect.
#


# Compile an expression and compare the result with a fixed value
test_compile() {
    input="$1"
    result="$2"

    #
    # Do this the long way round so we have assembly file for
    # inspection if/when a test fails.
    #
    go run main.go "${input}" > test.s
    gcc -static -o ./test test.s

    #
    # Run the test.
    #
    out=$(./test  | awk '{print $NF}')

    if [ "${result}" = "${out}" ]; then
        echo "Expected output found for '$input' [$result] "
        rm test test.s
    else
        echo "Expected output of '$input' is '$result' - got '${out}' instead"
        exit 1
    fi

}



# Simple operations
test_compile '3 4 +' 7
test_compile '3 4 *' 12
test_compile '10 2 -' 8
test_compile '10 2 /' 5

# modulus
#test_compile  '1 4 %' 1
#test_compile  '2 4 %' 2
#test_compile  '3 4 %' 3
#test_compile  '4 4 %' 0
#test_compile  '5 4 %' 1
#test_compile  '6 4 %' 2
#test_compile  '7 4 %' 3
#test_compile  '8 4 %' 0
#test_compile  '9 4 %' 1
#test_compile '10 4 %' 2
#test_compile '11 4 %' 3
#test_compile '12 4 %' 0

# powers of two - the manual-way
test_compile '2 2 *' 4
test_compile '2 2 * 2 *' 8
test_compile '2 2 * 2 * 2 *' 16
test_compile '2 2 * 2 * 2 * 2 *' 32
test_compile '2 2 * 2 * 2 * 2 * 2 *' 64
test_compile '2 2 * 2 * 2 * 2 * 2 * 2 *' 128
test_compile '2 2 * 2 * 2 * 2 * 2 * 2 * 2 *' 256
test_compile '2 2 * 2 * 2 * 2 * 2 * 2 * 2 * 2 *' 512
test_compile '2 2 * 2 * 2 * 2 * 2 * 2 * 2 * 2 * 2 *' 1024

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


# Note we're operating on integers, so these are "correct".
test_compile '3 2 /' 1.5
test_compile '5 2 /' 2.5


exit 0
