package streets

import (
	"errors"
	"github.com/rs/zerolog/log"
)

func (v *Vehicle) Drive() {
	for !v.IsParked {
		v.Step()
	}
	log.Info().Msgf("[%s] is parked.", v.ID)
}

func (v *Vehicle) Step() {
	edge, err := v.g.Graph.Edge(v.PrevID, v.NextID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get edge.")
		panic(err)
	}

	data, ok := edge.Properties.Data.(Data)
	if !ok {
		log.Error().Msg("Failed to convert edge data to Data.")
		panic(errors.New("failed to convert edge data to Data"))
	}

	v.DistanceRemaining = data.Length
	v.DistanceRemaining += v.Delta

	if v.DistanceRemaining >= v.Speed {
		v.DistanceRemaining -= v.Speed
	} else {
		v.Delta = v.DistanceRemaining
		v.DistanceRemaining = 0
	}

	v.PrevID = v.NextID
	v.NextID = v.getNextID()

	if v.NextID == 0 {
		v.IsParked = true
	}
}

// getNextID returns the next ID in the path, 0 if the vehicle is parked
func (v *Vehicle) getNextID() int {
	var prevIdIndex = -1

	for i := 0; i < len(v.PathIDs); i++ {
		if v.PathIDs[i] == v.PrevID {
			prevIdIndex = i
		}
	}

	if prevIdIndex == -1 {
		log.Error().Msgf("Could not find prevID %d in PathIDs %v", v.PrevID, v.PathIDs)
		panic(errors.New("could not find prevID in PathIDs"))
	}

	isLastIdx := prevIdIndex == len(v.PathIDs)-1

	// if the vehicle is parked, it is at its destination
	if v.NextID == 0 || isLastIdx || v.IsParked {
		// if vehicle is parked nextID is not 0
		return 0
	}

	nextID := v.PathIDs[prevIdIndex+1]

	vertexExistsOnCurrentGraph := v.g.VertexExists(nextID)
	if !vertexExistsOnCurrentGraph {
		// TODO: implement
		err := errors.New("failed to get vertex")
		log.Error().Err(err).Msgf("NOT IN GRAPH: Failed to get vertex %d", nextID)
		panic(err)
	}

	return nextID
}
