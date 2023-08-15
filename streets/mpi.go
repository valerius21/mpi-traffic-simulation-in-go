package streets

import (
	"errors"
	"github.com/rs/zerolog/log"
	mpi "github.com/sbromberger/gompi"
)

const (
	ROOT_ID              = 0
	VEHICLE_OUT_TAG      = 1
	VEHICLE_IN_ROOT_TAG  = 2
	VEHICLE_IN_LEAF_TAG  = 3
	REQUEST_EDGE         = 4
	RECEIVE_EDGE         = 5
	VEHICLE_OUT_ROOT_TAG = 6
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

// SendVehicleToRoot sends a vehicle to the root process using MPI Broadcast
func (m *MPI) SendVehicleToRoot(vehicle Vehicle) error {
	jBytes, err := vehicle.Marshal()
	if err != nil {
		return errors.New("failed to pack vehicle")
	}
	log.Info().Msgf("[%d] sending vehicle to root", m.taskID)
	m.comm.SendBytes(jBytes, ROOT_ID, VEHICLE_OUT_TAG)

	// TODO: delete vehicle from current graph
	return nil
}

func (m *MPI) EmitVehicle(vehicle Vehicle) error {
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	jBytes, err := vehicle.Marshal()
	if err != nil {
		return errors.New("failed to pack vehicle")
	}

	log.Info().Msgf("[%d] sending vehicle", m.taskID)

	// broadcast vehicle to all processes
	m.comm.BcastBytes(jBytes, ROOT_ID)

	log.Info().Msgf("[%d] sent vehicle", m.taskID)
	log.Warn().Msgf("[%d] vehicle deletion not implemented", m.taskID)

	return nil
}

func (m *MPI) ReceiveAndSendVehicleOverRoot(leafs []*StreetGraph) error {
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	jBytes, _ := m.comm.RecvBytes(mpi.AnySource, VEHICLE_OUT_TAG)

	vehicle, err := UnmarshalVehicle(jBytes)
	if err != nil {
		log.Error().Msgf("failed to unmarshal vehicle: %s", err.Error())
		return err
	}

	targetID, err := m.g.GetRectFromVertexID(vehicle.NextID, leafs)
	if err != nil {
		return err
	}
	if targetID == -1 {
		return errors.New("failed to find target leaf")
	}

	m.comm.SendBytes(jBytes, targetID, VEHICLE_IN_LEAF_TAG)
	return nil
}

func (m *MPI) ReceiveVehicleOnLeaf() (Vehicle, error) {
	jBytes, _ := m.comm.RecvBytes(ROOT_ID, VEHICLE_IN_LEAF_TAG)
	vehicle, err := UnmarshalVehicle(jBytes)
	if err != nil {
		return Vehicle{}, err
	}
	return vehicle, nil
}
