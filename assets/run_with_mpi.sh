#!/usr/bin/env bash

# install deps
# go mod tidy

# Run the program
mpirun -np 2 -quiet go run cmd/main.go -mpi "$@"
