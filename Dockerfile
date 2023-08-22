FROM ubuntu:23.04

# Update 
RUN apt update && apt install -y

# TODO: clean up unused packages
RUN apt install -y git build-essential bash zsh fish libopenmpi-dev openmpi-common htop golang wget

# Compile MPI
WORKDIR /tmp/mpi_src
RUN cd /tmp/mpi_src && wget https://download.open-mpi.org/release/open-mpi/v4.1/openmpi-4.1.5.tar.gz && tar xvf openmpi-4.1.5.tar.gz
RUN cd /tmp/mpi_src/openmpi-4.1.5 && ./configure --prefix=/usr/local && make all && make install

# WARN: Set it's okay with root
ENV OMPI_ALLOW_RUN_AS_ROOT=1
ENV OMPI_ALLOW_RUN_AS_ROOT_CONFIRM=1

# Get Go ready
WORKDIR /app

COPY . .
RUN go mod tidy
RUN go mod download
RUN go get golang.org/x/tools/cmd/stringer


# Keep it running
CMD ["tail", "-f", "/dev/null"]

