package streets

import (
	"bytes"
	"encoding/gob"
)

type LengthFloat struct {
	Length float64
}

func (l *LengthFloat) Marshal() ([]byte, error) {
	tmp := *l
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(tmp)
	return buf.Bytes(), err
}

func UnmarshalLengthFloat(data []byte) (LengthFloat, error) {
	var lengthFloat LengthFloat
	byteBuffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(byteBuffer)

	err := dec.Decode(&lengthFloat)
	return lengthFloat, err
}
