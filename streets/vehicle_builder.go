package streets

import (
	"errors"
	"github.com/aidarkhanov/nanoid"
	"github.com/rs/zerolog/log"
)

type VehicleBuilder struct {
	speed   float64
	pathIDs []int

	delta    float64
	isParked bool

	prevID int
	nextID int

	graph *StreetGraph
}

func NewVehicleBuilder() *VehicleBuilder {
	return &VehicleBuilder{}
}

func (vb *VehicleBuilder) WithSpeed(speed float64) *VehicleBuilder {
	vb.speed = speed
	return vb
}

func (vb *VehicleBuilder) WithPathIDs(pathIDs []int) *VehicleBuilder {
	vb.pathIDs = pathIDs
	return vb
}

func (vb *VehicleBuilder) WithGraph(graph *StreetGraph) *VehicleBuilder {
	vb.graph = graph
	return vb
}

func (vb *VehicleBuilder) WithLastID(lastID int) *VehicleBuilder {
	vb.prevID = lastID
	return vb
}

func (vb *VehicleBuilder) WithNextID(nextID int) *VehicleBuilder {
	vb.nextID = nextID
	return vb
}

func (vb *VehicleBuilder) WithDelta(delta float64) *VehicleBuilder {
	vb.delta = delta
	return vb
}

func (vb *VehicleBuilder) WithIsParked(isParked bool) *VehicleBuilder {
	vb.isParked = isParked
	return vb
}

func (vb *VehicleBuilder) FromJsonBytes(jsonBytes []byte) (*VehicleBuilder, error) {
	v, err := UnmarshalVehicle(jsonBytes)
	if err != nil {
		return &VehicleBuilder{}, nil
	}

	vb.speed = v.Speed
	vb.pathIDs = v.PathIDs
	vb.delta = v.Delta
	vb.isParked = v.IsParked
	vb.prevID = v.PrevID
	vb.nextID = v.NextID

	return vb, nil
}

func (vb *VehicleBuilder) check() (*VehicleBuilder, error) {
	if vb.speed == 0. {
		err := errors.New("speed is not set")
		log.Error().Err(err).Msg("Failed to build vehicle.")
		return nil, err
	}
	if len(vb.pathIDs) == 0 {
		err := errors.New("path is not set")
		log.Error().Err(err).Msg("Failed to build vehicle.")
		return nil, err
	}
	if vb.graph == nil {
		err := errors.New("graph is not set")
		log.Error().Err(err).Msg("Failed to build vehicle.")
		return nil, err
	}
	if vb.prevID < 1 {
		err := errors.New("prevID is not set")
		log.Error().Err(err).Msg("Failed to build vehicle.")
		return nil, err
	}
	return vb, nil
}

func (vb *VehicleBuilder) Build() (Vehicle, error) {
	vehicle := Vehicle{
		ID:                nanoid.New(),
		PathIDs:           vb.pathIDs,
		Speed:             vb.speed,
		Delta:             vb.delta,
		NextID:            vb.nextID,
		PrevID:            vb.prevID,
		IsParked:          vb.isParked,
		DistanceRemaining: 0.0, // default value
	}
	return vehicle, nil
}
