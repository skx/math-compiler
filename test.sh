#!/bin/sh
#
# Simple test-script to ensure we get the result we expect
#


test_compile() {
    input="$1"
    result="$2"

    go run main.go "${input}" > test.s
    gcc -static -o ./test test.s
    ./test

    if [ "${result}" = $? ]; then
        echo "Expected output found for '$input' $result "
    else
        echo "Expected output of '$input' is $result - got $?"
    fi

}


test_compile '3 4 +' 7
test_compile '3 4 *' 12
test_compile '10 2 -' 8
test_compile '10 2 /' 5
