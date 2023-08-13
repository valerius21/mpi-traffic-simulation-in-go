package streets

import "encoding/json"

func UnmarshalVehicle(data []byte) (Vehicle, error) {
	var r Vehicle
	err := json.Unmarshal(data, &r)
	return r, err
}

func (v *Vehicle) Marshal() ([]byte, error) {
	return json.Marshal(v)
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
	g                 *StreetGraph
}
