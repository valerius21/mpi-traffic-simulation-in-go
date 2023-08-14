package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mpi "github.com/sbromberger/gompi"
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

	if !mpi.IsOn() {
		if *useRoutines {
			runWithGoRoutines(vehicleList)
		} else {
			runSequentially(vehicleList)
		}
		return
	}

	mpi.Start(false)
	defer mpi.Stop()
	comm := mpi.NewCommunicator(nil)

	//numTasks := world.Size()
	taskID := comm.Rank()

	if mpi.WorldSize() < 2 {
		log.Error().Msg("World size is less than 2")
		return
	}

	rectangularSplits := mpi.WorldSize() - 1

	if taskID == 0 {
		size, err := rootGraph.Graph.Size()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get size of graph")
			return
		}
		log.Info().Msgf("Number of vertices: %d", size)
		m := streets.NewMPI(0, *comm, rootGraph)
		for {
			// root process will listen for incoming requests
			err = m.RespondToEdgeLengthRequest()
			if err != nil {
				log.Error().Err(err).Msg("Failed to respond to edge length request")
				return
			}
		}
	} else {
		// Broadcast to the leafs, if one of them has the current vertex of the vehicle
		// just send it -> (P) What if the vertex does not exist? -> Vehicle is lost...
		// leafs will listen for incoming vehicles and add them to their graph
		// leafs start the drive on them
		// While true loop for constantly listening to receive bytes?

		_ = streets.NewMPI(taskID, *comm, rootGraph)

		_, ok := setupLeaf(jsonPath, rootGraph, rectangularSplits, taskID-1, taskID)
		if !ok {
			log.Error().Msgf("[%d] Failed to setup leaf", taskID)
			return
		}
	}
}

func setupLeaf(jsonPath *string, rootGraph *streets.StreetGraph, rectangularSplits int, i int, taskID int) (*streets.StreetGraph, bool) {
	gb := streets.NewGraphBuilder().FromJsonFile(*jsonPath).IsLeaf(rootGraph).NumberOfRects(rectangularSplits)
	gb = gb.PickRect(i - 1).DivideGraphsIntoRects().FilterForRect()
	gb = gb.SetTopRightBottomLeftVertices()
	leafGraph, err := gb.Build()
	if err != nil {
		log.Error().Err(err).Msg("Failed to build graph")
		return nil, true
	}
	size, err := leafGraph.Graph.Size()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get size of graph")
		return nil, true
	}
	log.Info().Msgf("[%d] Number of vertices: %d", taskID, size)
	return leafGraph, false
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
