#!/usr/bin/env bash

# install deps
# go mod tidy

# Run the program
mpirun -np 3 go run cmd/main.go -mpi "$@"
