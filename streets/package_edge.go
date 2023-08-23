package streets

import (
	"bytes"
	"encoding/gob"
)

func UnmarshalEdgePackage(data []byte) (EdgePackage, error) {
	var r EdgePackage
	byteBuffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(byteBuffer)

	err := dec.Decode(&r)

	return r, err
}

func (r *EdgePackage) Marshal() ([]byte, error) {
	tmp := *r

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(tmp)
	return buf.Bytes(), err
}

type EdgePackage struct {
	Src  int `json:"src"`
	Dest int `json:"dest"`
}
