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
	REQUEST_DONE_INC_TAG = 7
	DONE_BCAST_TAG       = 8
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

	edgePackage, err := e.Marshal()
	if err != nil {
		return 0, errors.New("failed to pack edge package")
	}

	log.Info().Msgf("[%d] sending edge package len(%d)", m.taskID, len(edgePackage))
	// send request to root
	log.Debug().Msgf("[%d] sending edge %d->%d to root", m.taskID, e.Src, e.Dest)
	m.comm.SendBytes(edgePackage, ROOT_ID, REQUEST_EDGE)

	log.Info().Msgf("[%d] waiting to get edge package from root", m.taskID)
	// receive edge length from root
	//TODO: length, _ := m.comm.RecvFloat64(ROOT_ID, RECEIVE_EDGE)
	bytes, status := m.comm.RecvBytes(ROOT_ID, RECEIVE_EDGE)
	log.Info().Msgf("[%d] received edge package from %d", m.taskID, status.GetSource())
	lf, err := UnmarshalLengthFloat(bytes)
	if err != nil {
		log.Error().Msgf("failed to unmarshal length float: %s", err.Error())
		return 0, errors.New("failed to unmarshal length float")
	}

	length := lf.Length

	if length <= 0.0 {
		return 0, errors.New("failed to pack edge package")
	}

	return length, nil
}

func (m *MPI) RespondToEdgeLengthRequest() error {
	log.Info().Msgf("[%d] waiting for edge package - %d", m.taskID)
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	log.Info().Msg("[root] waiting for edge package")

	bytes, status := m.comm.RecvBytes(2, REQUEST_EDGE) // FIXME: MPI.ANY_SOURCE
	edgePackage, err := UnmarshalEdgePackage(bytes)

	log.Info().Msgf("[root] received edge package from %d", status.GetSource())

	if err != nil {
		return errors.New("failed to unmarshal edge package")
	}

	log.Debug().Msgf("[root] received edge package from %d src(%d) dest(%d)", status.GetSource(), edgePackage.Src,
		edgePackage.Dest)
	edge, err := m.g.Graph.Edge(edgePackage.Src, edgePackage.Dest)

	if err != nil {
		log.Error().Msgf("failed to get edge: %s", err.Error())
		return err
	}

	data, ok := edge.Properties.Data.(Data)
	if !ok {
		return errors.New("edge data is not of type Data")
	}

	log.Info().Msgf("[root] sending edge package %f", data.Length)

	// send edge length to sender
	lf := LengthFloat{Length: data.Length}
	lfBytes, err := lf.Marshal()
	if err != nil {
		log.Error().Msgf("failed to marshal length float: %s", err.Error())
		return errors.New("failed to pack length float")
	}
	m.comm.SendBytes(lfBytes, status.GetSource(), RECEIVE_EDGE)

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

func (m *MPI) EmitVehicle(vehicle Vehicle, lookupTable map[int]int) error {
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	jBytes, err := vehicle.Marshal()
	if err != nil {
		return errors.New("failed to pack vehicle")
	}

	log.Debug().Msgf("[%d] sending vehicle", m.taskID)

	// broadcast vehicle to all processes
	targetID := lookupTable[vehicle.NextID]

	if targetID == 0 {
		return errors.New("failed to find target leaf")
	}

	m.comm.SendBytes(jBytes, targetID, VEHICLE_IN_LEAF_TAG)

	log.Info().Msgf("[%d] sent vehicle - %s", m.taskID, vehicle.ID)

	return nil
}

func (m *MPI) ReceiveAndSendVehicleOverRoot(lookupTable map[int]int) error {
	if m.taskID != ROOT_ID {
		return errors.New("process is not root")
	}

	jBytes, status := m.comm.RecvBytes(mpi.AnySource, VEHICLE_OUT_TAG)
	log.Debug().Msgf("[%d] received vehicle request from %d", m.taskID, status.GetSource())
	vehicle, err := UnmarshalVehicle(jBytes)
	log.Info().Msgf("[%d] received vehicle from %d", m.taskID, status.GetSource())

	if err != nil {
		log.Error().Msgf("failed to unmarshal vehicle: %s", err.Error())
		return err
	}
	targetID := lookupTable[vehicle.NextID]
	//targetID, err := m.g.GetRectFromVertexID(vehicle.NextID, leafs)
	if targetID == 0 {
		return errors.New("failed to find target leaf")
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

func (m *MPI) SendDoneToRoot() {
	m.comm.SendInt32(int32(1), ROOT_ID, REQUEST_DONE_INC_TAG)
}

func (m *MPI) ReceiveDoneFromLeaf(incrementor *int) {
	log.Warn().Msgf("[%d] waiting for done from leaf", m.taskID)
	b, _ := m.comm.RecvInt32(mpi.AnySource, REQUEST_DONE_INC_TAG)
	log.Warn().Msgf("[%d] received done from leaf -> %v", m.taskID, b)
	if b == 1 {
		v := *incrementor
		v++
		*incrementor = v
	}
}

func (m *MPI) BCastDone() int32 {
	doneArr := make([]int32, 1)
	if m.taskID == ROOT_ID {
		doneArr[0] = int32(1)
	}

	m.comm.BcastInt32s(doneArr, ROOT_ID)
	return 1
}
