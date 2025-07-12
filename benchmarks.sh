#!/bin/bash

set -e

# Execute the benchmarks
go test -bench=. -run=^$ ./...
