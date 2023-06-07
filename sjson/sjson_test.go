package sjson

import (
	"fmt"
	"testing"

	json "github.com/bytedance/sonic"
)

func BenchmarkSet(b *testing.B) {
	json := []byte(`{"name":{"first":"Janet","last":"Prichard"},"age":47}`)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		SetBytes(json, "gender", "Male")
	}
}

func BenchmarkJsonMarshalUnmarshal(b *testing.B) {
	const js = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		var tb map[string]any
		json.Unmarshal([]byte(js), &tb)
		tb["gender"] = "Male"
		json.Marshal(tb)
	}
}

func BenchmarkArray(b *testing.B) {
	json := `{
    "app": "gfgeeks",
    "prop": [
          {"region": 736,"set": true,"score": 72},
          {"region": 563,"set": true,"score": 333},
          {"region": 563,"set": false,"score": 333}
    ],
    "index" : "haskell"
}`
	for n := 0; n < b.N; n++ {
		if val := Get(json, "prop"); val.Exists() && val.IsArray() {
			for i, _ := range val.Array() {
				json, _ = Set(json, fmt.Sprintf("prop.%d.name", i), "name")
			}
		}
	}
}

func BenchmarkJsonMap(b *testing.B) {
	js := `{
    "app": "gfgeeks",
    "prop": [
          {"region": 736,"set": true,"score": 72},
          {"region": 563,"set": true,"score": 333},
          {"region": 563,"set": false,"score": 333}
    ],
    "index" : "haskell"
}`
	for n := 0; n < b.N; n++ {
		var tb map[string]any
		json.Unmarshal([]byte(js), &tb)
		if v, ok := tb["prop"]; ok {
			l := v.([]any)
			for _, j := range l {
				switch j := j.(type) {
				case map[string]any:
					j["name"] = "name"
				}

			}
		}
		json.Marshal(tb)
	}
}
