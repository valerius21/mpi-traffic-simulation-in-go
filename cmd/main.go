package main

import (
	"encoding/json"
	"fmt"
)

type MySpeed struct {
	Speed float64 `json:"speed"`
}

func main() {
	data := `{}`
	var mySpeed MySpeed
	err := json.Unmarshal([]byte(data), &mySpeed)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(mySpeed.Speed) // This will print "0"
}
