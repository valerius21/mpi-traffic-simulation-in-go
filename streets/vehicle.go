package streets

import (
	"errors"
	"github.com/rs/zerolog/log"
)

// Drive is for the non-MPI implementation
func (v *Vehicle) Drive() {
	for !v.IsParked {
		v.Step()
	}
	log.Info().Msgf("[%s] is parked.", v.ID)
}

// Step is the main algorithm for the vehicle
func (v *Vehicle) Step() {
	log.Debug().Msgf("[%s] is stepping.", v.ID)
	if v.NextID < 0 { // III.1
		log.Debug().Msgf("[%s] is parked. (III.3)", v.ID)
		v.IsParked = true
	}

	// III.2 Assuming that the edge is in the graph
	edge, err := v.StreetGraph.Graph.Edge(v.PrevID, v.NextID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get edge.")
		panic(err)
	}
	log.Debug().Msgf("[%s] is on edge %d -> %d (III.2)", v.ID, v.PrevID, v.NextID)

	data, ok := edge.Properties.Data.(Data)
	if !ok {
		log.Error().Msg("Failed to convert edge data to Data.")
		panic(errors.New("failed to convert edge data to Data"))
	}

	v.DistanceRemaining = data.Length

	// III.3
	v.DistanceRemaining += v.Delta
	log.Debug().Msgf("[%s] has distance remaining %f (III.3)", v.ID, v.DistanceRemaining)

	// III.4
	log.Debug().Msgf("[%s] has speed %f (III.4)", v.ID, v.Speed)
	for v.DistanceRemaining >= v.Speed && v.DistanceRemaining-v.Speed > 0 {
		v.DistanceRemaining -= v.Speed // III.5
		log.Debug().Msgf("[%s] has distance remaining %f (III.5)", v.ID, v.DistanceRemaining)
	}
	// III.6
	v.Delta = v.DistanceRemaining
	v.DistanceRemaining = 0
	log.Debug().Msgf("[%s] has delta remaining %f (III.6)", v.ID, v.Delta)

	// because no vertex ID can be -1, which indicates a leaf switch.
	nextStepId := v.GetNextID(v.NextID)
	if nextStepId == -1 {
		// III.9.2
		log.Info().Msgf("[%s] is marked for deletion. (III.9.2)", v.ID)
		return
	} else if nextStepId == 0 {
		// III.8
		log.Info().Msgf("[%s] is parked. (III.8)", v.ID)
		v.IsParked = true
		return
	}

	v.PrevID = v.NextID // III.6.1
	log.Debug().Msgf("[%s] has prevID %d (III.6.1)", v.ID, v.PrevID)
	v.NextID = v.GetNextID(v.PrevID) // III.6.2
	log.Debug().Msgf("[%s] has nextID %d (III.6.2)", v.ID, v.NextID)

	// III.8
	if v.NextID == 0 {
		log.Debug().Msgf("[%s] is parked. (III.8)", v.ID)
		v.IsParked = true
		return
	}

	// III.9.1
	// continue steps
	log.Debug().Msgf("[%s] is continuing steps. (III.9.1)", v.ID)
}

// GetNextID returns the next ID in the path, 0 if the vehicle is parked (III.7)
func (v *Vehicle) GetNextID(prevID int) int {
	var prevIdIndex = -1

	for i := 0; i < len(v.PathIDs); i++ {
		if v.PathIDs[i] == prevID {
			prevIdIndex = i
		}
	}

	isLastIdx := prevIdIndex == len(v.PathIDs)-1

	// if the vehicle is parked, it is at its destination
	if v.NextID == 0 || isLastIdx || v.IsParked {
		// if vehicle is parked nextID is not 0
		v.IsParked = true
		return 0
	}

	nextID := v.PathIDs[prevIdIndex+1]

	vertexExistsOnCurrentGraph := v.StreetGraph.VertexExists(nextID)
	if !vertexExistsOnCurrentGraph {
		// III.9.2
		v.MarkedForDeletion = true
		return -1
	}

	return nextID
}
