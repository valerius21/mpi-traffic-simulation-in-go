# PCHPC Project - Traffic Simulation with Go

### Development

Build the running container:

```bash
docker build . -t vmpi:dev
```

Running the Simulation:

```bash
docker run -d --name vmpi -v $(pwd):/app vmpi:dev
```

```bash
# run 10 vehicles
docker exec -it vmpi /app/assets/run_with_mpi.sh -n 10 -debug
```

# Archived
The working rewrite in Rust can be found here: https://github.com/valerius21/mpi-traffic-sim-rust/
