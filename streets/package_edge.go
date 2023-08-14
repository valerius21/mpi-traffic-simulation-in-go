package streets

import (
	"encoding/json"
)

func UnmarshalEdgePackage(data []byte) (EdgePackage, error) {
	var r EdgePackage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *EdgePackage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type EdgePackage struct {
	Src  int `json:"src"`
	Dest int `json:"dest"`
}

func (r *EdgePackage) Pack() ([]byte, error) {
	return r.Marshal()
}
