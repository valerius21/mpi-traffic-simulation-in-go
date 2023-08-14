#!/usr/bin/env bash

# install deps
go mod tidy

# Run the program
mpirun -np 4 go test --count=1 ./...
