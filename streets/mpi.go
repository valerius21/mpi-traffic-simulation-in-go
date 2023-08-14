package streets

import (
	"errors"
	"github.com/rs/zerolog/log"
	mpi "github.com/sbromberger/gompi"
)

const (
	IsMPI           = false
	ROOT_ID         = 0
	VEHICLE_OUT_TAG = 1
	VEHICLE_IN_TAG  = 2
	REQUEST_EDGE    = 3
	RECEIVE_EDGE    = 4
)

type MPI struct {
	taskID int
	comm   mpi.Communicator
	g      *StreetGraph
}

func NewMPI(taskID int, communicator mpi.Communicator, graph *StreetGraph) *MPI {
	return &MPI{taskID: taskID, comm: communicator, g: graph}
}

func (m *MPI) AskRootForEdgeLength(srcVertexID, destVertexID int) (float64, error) {
	if m.taskID == ROOT_ID {
		// process is root
		return 0, errors.New("process is root")
	}

	// package
	e := EdgePackage{
		Src:  srcVertexID,
		Dest: destVertexID,
	}

	edgePackage, err := e.Pack()
	if err != nil {
		return 0, errors.New("failed to pack edge package")
	}

	log.Info().Msgf("[%d] sending edge package", m.taskID)
	// send request to root
	m.comm.SendBytes(edgePackage, ROOT_ID, REQUEST_EDGE)

	log.Info().Msgf("[%d] waiting to get edge package from root", m.taskID)
	// receive edge length from root
	length, _ := m.comm.RecvFloat64(ROOT_ID, RECEIVE_EDGE)

	if length <= 0.0 {
		return 0, errors.New("failed to pack edge package")
	}

	return length, nil
}

func (m *MPI) RespondToEdgeLengthRequest() error {
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	log.Info().Msg("[root] waiting for edge package")
	jBytes, status := m.comm.RecvBytes(mpi.AnySource, mpi.AnyTag)
	log.Info().Msgf("[root] received edge package from %d", status.GetSource())
	if jBytes == nil {
		return errors.New("failed to receive edge package")
	}

	edgePackage, err := UnmarshalEdgePackage(jBytes)
	if err != nil {
		return errors.New("failed to unmarshal edge package")
	}

	edge, err := m.g.Graph.Edge(edgePackage.Src, edgePackage.Dest)

	if err != nil {
		return err
	}

	data, ok := edge.Properties.Data.(Data)
	if !ok {
		return errors.New("edge data is not of type Data")
	}

	log.Info().Msg("[root] sending edge package")
	// send edge length to sender
	m.comm.SendFloat64(data.Length, status.GetSource(), RECEIVE_EDGE)

	return nil
}
