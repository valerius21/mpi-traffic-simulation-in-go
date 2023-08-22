package streets

import (
	"errors"
	"github.com/aidarkhanov/nanoid"
	"github.com/rs/zerolog/log"
	"strings"
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
	if len(vb.pathIDs) < 2 {
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
	vb, err := vb.check()

	if err != nil {
		return Vehicle{}, err
	}

	newAlphabet := nanoid.DefaultAlphabet
	newAlphabet = strings.Replace(newAlphabet, "_", "", -1)
	newAlphabet = strings.Replace(newAlphabet, "-", "", -1)
	vid, err := nanoid.Generate(newAlphabet, 10)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate vehicle ID.")
		return Vehicle{}, err
	}

	vehicle := Vehicle{
		ID:                vid,
		PathIDs:           vb.pathIDs,
		Speed:             vb.speed,
		Delta:             vb.delta,
		NextID:            vb.nextID,
		PrevID:            vb.prevID,
		IsParked:          vb.isParked,
		DistanceRemaining: 0.0, // default value
		StreetGraph:       vb.graph,
	}

	// ensure nextID is set
	id := vehicle.GetNextID(vehicle.PathIDs[0])
	if id == 0 {
		vehicle.IsParked = true
		vehicle.NextID = 0
	} else if id == -1 {
		// TODO: this should never happen
		panic("vehicle.GetNextID returned -1 at initialization")
	} else {
		vehicle.NextID = id
	}

	return vehicle, nil
}
