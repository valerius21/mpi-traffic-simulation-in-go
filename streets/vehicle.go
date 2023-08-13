package streets

//
//import (
//	"errors"
//	"fmt"
//	"github.com/aidarkhanov/nanoid"
//	"github.com/dominikbraun/graph"
//	"github.com/rs/zerolog/log"
//	"math"
//	"pchpc_next/utils"
//	"sort"
//)
//
//// OVehicle is a vehicle
//type OVehicle struct {
//	ID                string
//	Path              []int
//	DistanceTravelled float64
//	Speed             float64
//	g                 *StreetGraph
//	IsParked          bool
//	PathLengths       []float64
//	PathLimit         float64
//}
//
//// getPathLengths calculates the length of each edge in the path
//func (v *OVehicle) getPathLengths() error {
//	lengthsArray := make([]float64, 0)
//	sum := 0.0
//
//	for i, vertex := range v.Path {
//		if i == len(v.Path)-1 {
//			break
//		}
//
//		edge, err := (*v.g).Edge(vertex, v.Path[i+1])
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get edge.")
//			return err
//		}
//
//		data, ok := edge.Properties.Data.(Data)
//		if !ok {
//			log.Error().Msg("Failed to convert edge data to Data.")
//			return err
//		}
//
//		lengthsArray = append(lengthsArray, data.Length)
//		sum += data.Length
//	}
//
//	v.PathLengths = lengthsArray
//	v.PathLimit = sum
//	return nil
//}
//
//// getCurrentEdge returns the current edge the vehicle is on
//func (v *OVehicle) getCurrentEdge() (*graph.Edge[JVertex], error) {
//	idx, _ := v.deductCurrentPathVertexIndex()
//	edge, err := v.getEdgeByIndex(idx)
//	if err != nil {
//		return nil, err
//	}
//	return edge, nil
//}
//
//// deductCurrentPathVertexIndex returns the index of the current edge in the path, delta must be added
//func (v *OVehicle) deductCurrentPathVertexIndex() (index int, delta float64) {
//	tmpDistance := v.DistanceTravelled
//
//	for i, length := range v.PathLengths {
//		if tmpDistance < length {
//			return i, math.Abs(tmpDistance)
//		}
//		tmpDistance -= length
//	}
//
//	return 0, 0.0
//}
//
//// getEdgeByIndex returns the edge at the given index
//func (v *OVehicle) getEdgeByIndex(index int) (oEdge *graph.Edge[JVertex], err error) {
//	if index == len(v.Path)-1 {
//		return oEdge, fmt.Errorf("index is out of range")
//	}
//
//	ed, err := (*v.g).Edge(v.Path[index], v.Path[index+1])
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to get edge.")
//		return &ed, err
//	}
//
//	return &ed, nil
//}
//
//// getHashMapByEdge returns the hashmap of the given edge
//func (v *OVehicle) getHashMapByEdge(edge *graph.Edge[JVertex]) (*utils.HashMap[string, *OVehicle], error) {
//	data, exists := edge.Properties.Data.(Data)
//	if !exists {
//		err := fmt.Errorf("edge data is not of type EdgeData")
//		log.Error().Err(err).Msg("Failed to get data from edge.")
//		return nil, err
//	}
//	return data.Map, nil
//}
//
//// isInMap checks if the vehicle is in the given hashmap
//func (v *OVehicle) isInMap(hashMap *utils.HashMap[string, *OVehicle]) bool {
//	_, exists := hashMap.Get(v.ID)
//	return exists
//}
//
//// AddVehicleToEdge adds the vehicle to the given hashmap
//func (v *OVehicle) AddVehicleToEdge(edge *graph.Edge[JVertex]) error {
//	edgeData, ok := edge.Properties.Data.(Data)
//	if !ok {
//		err := fmt.Errorf("edge data is not of type EdgeData")
//		log.Error().Err(err).Msg("Failed to get data from edge.")
//		return err
//	}
//
//	msEdgeSpeed := edgeData.MaxSpeed / 3.6
//	hashMap := edgeData.Map
//	if v.isInMap(hashMap) {
//		return nil
//	}
//
//	frontVehicle, err := v.GetFrontVehicleFromEdge(edge)
//	if err != nil {
//		return err
//	}
//
//	if frontVehicle != nil && frontVehicle.Speed < v.Speed {
//		v.Speed = frontVehicle.Speed
//	} else if frontVehicle != nil && frontVehicle.Speed > v.Speed && msEdgeSpeed > v.Speed {
//		minAcceleration := 0.1
//		maxAcceleration := 0.5
//		v.Speed += utils.RandomFloat64(minAcceleration, maxAcceleration)
//	}
//
//	hashMap.Set(v.ID, v)
//	v.updateVehiclePosition(hashMap)
//	return nil
//}
//
//// RemoveVehicleFromMap removes the vehicle from the given hashmap
//func (v *OVehicle) RemoveVehicleFromMap(hashMap *utils.HashMap[string, *OVehicle]) {
//	if hashMap.Len() == 0 {
//		v.updateVehiclePosition(hashMap)
//		return
//	}
//	hashMap.Del(v.ID)
//	v.updateVehiclePosition(hashMap)
//}
//
//// updateVehiclePosition updates the vehicle position
//func (v *OVehicle) updateVehiclePosition(hashMap *utils.HashMap[string, *OVehicle]) {
//	if v.PathLimit <= v.DistanceTravelled {
//		v.IsParked = true
//		edge, err := v.getCurrentEdge()
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get current edge.")
//			return
//		}
//		hashMap, err := v.getHashMapByEdge(edge)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get hashmap.")
//			return
//		}
//
//		hashMap.Del(v.ID)
//	}
//
//	log.Debug().Msgf("Current vehicles on edge: %d, %s", hashMap.Len(), v.ID)
//}
//
//// String returns the string representation of the vehicle
//func (v *OVehicle) String() string {
//	return fmt.Sprintf("OVehicle: %s, Speed: %f, Distance Travelled: %v Sum: %.2f", v.ID, v.Speed,
//		v.DistanceTravelled, v.PathLimit)
//}
//
//// Step moves the vehicle one step forward
//func (v *OVehicle) Step() {
//	idx, delta := v.deductCurrentPathVertexIndex()
//	log.Debug().Msgf("Current index: %d, delta: %f", idx, delta)
//	log.Debug().Msgf("Current path: %v", v.Path)
//	log.Debug().Msgf("Current vehicle: %v", v)
//	edge, err := v.getEdgeByIndex(idx)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to get edge.")
//		return
//	}
//
//	if v.Speed >= delta && idx != 0 {
//		// check if vertex exists
//		_, err := (*v.g).Vertex(edge.Target.ID)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get vertex.")
//			panic(err)
//		}
//
//		oldEdge, err := v.getEdgeByIndex(idx - 1)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get edge.")
//			return
//		}
//
//		oldHashMap, err := v.getHashMapByEdge(oldEdge)
//		if err != nil {
//			log.Error().Err(err).Msg("Failed to get hashmap.")
//			return
//		}
//		v.RemoveVehicleFromMap(oldHashMap)
//	}
//
//	hashMap, err := v.getHashMapByEdge(edge)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to get hashmap.")
//		return
//	}
//	err = v.AddVehicleToEdge(edge)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to add vehicle to map.")
//		return
//	}
//	// vehicle is at destination
//	if v.IsParked {
//		return
//	}
//	v.drive()
//	v.updateVehiclePosition(hashMap)
//}
//
//// drive moves the vehicle forward
//func (v *OVehicle) drive() {
//	v.DistanceTravelled += v.Speed
//}
//
//// PrintInfo prints the vehicle info
//func (v *OVehicle) PrintInfo() {
//	log.Debug().
//		Str("id", v.ID).
//		Bool("isParked", v.IsParked).
//		Float64("speed", v.Speed).
//		Str("path lengths", fmt.Sprintf("%v", v.PathLengths)).
//		Msg("OVehicle info")
//}
//
//// GetFrontVehicleFromEdge returns the vehicle in front of the given vehicle
//func (v *OVehicle) GetFrontVehicleFromEdge(edge *graph.Edge[JVertex]) (*OVehicle, error) {
//	edgeData, ok := edge.Properties.Data.(Data)
//
//	if !ok {
//		return nil, errors.New("failed to cast edge data to Data")
//	}
//
//	eMap := edgeData.Map
//
//	if eMap.Len() < 1 {
//		return nil, nil
//	}
//
//	lst := eMap.ToList()
//
//	sort.Slice(lst, func(i, j int) bool {
//		return lst[i].DistanceTravelled > lst[j].DistanceTravelled
//	})
//
//	var frontIndex int
//
//	for i, vh := range lst {
//		if v.ID == vh.ID && i < 0 {
//			frontIndex = i - 1
//		}
//	}
//
//	return lst[frontIndex], nil
//}
//
//// NewVehicle creates a new vehicle
//func NewVehicle(speed float64, path []int, graph *graph.Graph[int, JVertex]) OVehicle {
//	v := OVehicle{
//		ID:                nanoid.New(),
//		Path:              path,
//		Speed:             speed,
//		g:                 graph,
//		DistanceTravelled: 0.0,
//	}
//	err := v.getPathLengths()
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to get path lengths.")
//		return OVehicle{}
//	}
//
//	return v
//}
