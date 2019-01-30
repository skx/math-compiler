#!/bin/sh
#
# Simple test-script to ensure we get the result we expect
#


test_compile() {
    input="$1"
    result="$2"

    go run main.go "${input}" > test.s
    gcc -static -o ./test test.s
    out=$(./test  | awk '{print $NF}')

    if [ "${result}" = "${out}" ]; then
        echo "Expected output found for '$input' [$result] "
        rm test test.s
    else
        echo "Expected output of '$input' is '$result' - got '${out}' instead"
    fi

}


test_compile '3 4 +' 7
test_compile '3 4 *' 12
test_compile '10 2 -' 8
test_compile '10 2 /' 5
test_compile '16384 2 * ' 32768
test_compile '16384 2 * 2 *' 65536
test_compile '16384 2 / 2 /' 4096

# We're operating on integers...
test_compile '3 2 /' 1
test_compile '5 2 /' 2
