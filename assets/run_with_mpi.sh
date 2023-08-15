#!/usr/bin/env bash

# install deps
go mod tidy

# Run the program
mpirun -np 4 go run cmd/main.go -mpi -debug