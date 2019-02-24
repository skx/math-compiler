#!/bin/sh


pwd
ls -l
# Run golang tests
go test ./...

# Run functional test-cases
./test.sh
