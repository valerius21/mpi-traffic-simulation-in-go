package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"pchpc_next/streets"
	"strconv"
	"sync"
)

func main() {
	// Flags
	n := flag.Int("n", 100, "Number of vehicles")
	useRoutines := flag.Bool("m", false, "Use goroutines")
	minSpeed := flag.Float64("min-speed", 5.5, "Minimum speed")
	maxSpeed := flag.Float64("max-speed", 8.5, "Maximum speed")
	jsonPath := flag.String("jsonPath", "assets/out.json", "Path to the json containing the graph data")
	debug := flag.Bool("debug", false, "Enable debug mode")
	numberOfRects := flag.Int("rects", 4, "Number of rects to divide the graph into")
	//useMPI := flag.Bool("mpi", false, "Use M3.1415...")

	flag.Parse()

	setupLogging(debug)

	b := streets.NewGraphBuilder().FromJsonFile(*jsonPath).SetTopRightBottomLeftVertices()
	rootGraph, err := b.NumberOfRects(1).DivideGraphsIntoRects().PickRect(0).IsRoot().Build()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build graph")
		return
	}

	// Create vehicles and drive
	ns := strconv.Itoa(*n)
	log.Info().Msg("Starting vehicles " + ns)

	vehicleList := make([]*streets.Vehicle, *n)

	if connectVehiclesToGraph(n, rootGraph, minSpeed, maxSpeed, vehicleList) {
		log.Error().Err(err).Msg("Failed to add vehicle")
		return
	}

	if *useRoutines {
		runWithGoRoutines(vehicleList)
	} else {
		runSequentially(vehicleList)
	}
}

func runWithGoRoutines(vehicleList []*streets.Vehicle) {
	var wg sync.WaitGroup
	for _, vehicle := range vehicleList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, vehicle *streets.Vehicle) {
			vehicle.Drive()
			wg.Done()
		}(&wg, vehicle)
	}
	wg.Wait()
}

func runSequentially(vehicleList []*streets.Vehicle) {
	for _, vehicle := range vehicleList {
		vehicle.Drive()
	}
}

func connectVehiclesToGraph(n *int, rootGraph *streets.StreetGraph, minSpeed *float64, maxSpeed *float64, vehicleList []*streets.Vehicle) bool {
	for i := 0; i < *n; i++ {
		v, err := rootGraph.AddVehicle(*minSpeed, *maxSpeed)
		if err != nil {
			return true
		}
		vehicleList[i] = v
	}
	return false
}

func setupLogging(debug *bool) {
	// Logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	runLogFile, _ := os.OpenFile(
		"main.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o664,
	)
	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
