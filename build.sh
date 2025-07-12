#!/bin/bash
set -e

# Execute the tests
go test ./...

# Evaluate the test results
# If the tests fail, exit with an error code
if [ $? -ne 0 ]; then
    echo "Tests failed. Exiting build process."
    exit 1
fi

# Build main package (CLI/game)
go build -o goengine ./cmd/main.go

# Build shared library from export/export.go
go build -buildmode=c-shared -o libgoengine.so ./export/export.go
