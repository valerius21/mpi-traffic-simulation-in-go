package streets

import "encoding/json"

func UnmarshalVehicle(data []byte) (Vehicle, error) {
	var r rawVehicle
	err := json.Unmarshal(data, &r)
	return Vehicle{
		ID:                r.ID,
		PathIDs:           r.PathIDs,
		Speed:             r.Speed,
		Delta:             r.Delta,
		NextID:            r.NextID,
		PrevID:            r.PrevID,
		IsParked:          r.IsParked,
		DistanceRemaining: r.DistanceRemaining,
		StreetGraph:       nil,
		MarkedForDeletion: false,
	}, err
}

func (v *Vehicle) Marshal() ([]byte, error) {
	return json.Marshal(rawVehicle{
		ID:                v.ID,
		PathIDs:           v.PathIDs,
		Speed:             v.Speed,
		Delta:             v.Delta,
		NextID:            v.NextID,
		PrevID:            v.PrevID,
		IsParked:          v.IsParked,
		DistanceRemaining: v.DistanceRemaining,
	})
}

type rawVehicle struct {
	ID                string  `json:"id"`
	PathIDs           []int   `json:"path_ids"`
	Speed             float64 `json:"speed"`
	Delta             float64 `json:"delta"`
	NextID            int     `json:"next_id"`
	PrevID            int     `json:"prev_id"`
	IsParked          bool    `json:"is_parked"`
	DistanceRemaining float64 `json:"distance_remaining"`
}

type Vehicle struct {
	ID                string  `json:"id"`
	PathIDs           []int   `json:"path_ids"`
	Speed             float64 `json:"speed"`
	Delta             float64 `json:"delta"`
	NextID            int     `json:"next_id"`
	PrevID            int     `json:"prev_id"`
	IsParked          bool    `json:"is_parked"`
	DistanceRemaining float64 `json:"distance_remaining"`
	StreetGraph       *StreetGraph
	MarkedForDeletion bool
}
