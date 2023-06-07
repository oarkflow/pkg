package main

import (
	"fmt"

	"github.com/oarkflow/pkg/sjson"
)

func main() {
	update()

}

func update() {
	json := `
{
    "app": "gfgeeks",
    "prop": [
          {"region": 736,"set": true,"score": 72},
          {"region": 563,"set": true,"score": 333},
          {"region": 563,"set": false,"score": 333}
    ],
    "index" : "haskell"
}`

	// loop through the "prop" values and find the target
	var index int
	var found bool
	if val := sjson.Get(json, "prop"); val.Exists() && val.IsArray() {
		for i, _ := range val.Array() {
			json, _ = sjson.Set(json, fmt.Sprintf("prop.%d.name", i), "name")
		}
	}
	fmt.Println(json)
	sjson.Get(json, "prop").ForEach(func(i, value sjson.Result) bool {
		json, _ = sjson.Set(json, fmt.Sprintf("prop.%d.name", i.Int()), "name")
		if value.Get("region").Int() == 563 && value.Get("set").Bool() {
			found = true
			return false
		}
		index++
		return true
	})
	if found {
		// if found the use sjson to update the value at index
		json, _ = sjson.Set(json, fmt.Sprintf("prop.%d.score", index), 334)
	}
	println(json)
}
