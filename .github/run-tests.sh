#!/bin/bassh

# Init the modules
go mod init

# Install the lint-tool
go get -u golang.org/x/lint/golint

# At this point failures cause aborts
set -e

# Run the linter
golint -set_exit_status ./...

# Run our golang tests
go test ./...

# Run functional test-cases
./test.sh
