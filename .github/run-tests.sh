#!/bin/sh

# init modules
go mod init

# Run golang tests
go test ./...

# Run functional test-cases
./test.sh
