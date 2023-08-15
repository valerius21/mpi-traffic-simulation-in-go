#!/usr/bin/env bash

# install deps
go mod tidy

# Run the program
mpirun -np 2 go run cmd/main.go -debug -mpi