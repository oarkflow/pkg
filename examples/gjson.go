package main

import (
	"encoding/json"
	"fmt"

	"github.com/oarkflow/pkg/sjson"
)

func main() {
	// update()
	test()
}

var requestData2 = []byte(`
{
    "patient": {
        "dob": "2003-04-10"
    },
  "coding": [
    {
      "em": {
          "code": "123",
          "encounter_uid": 1,
          "work_item_uid": 2, 
          "billing_provider": "Test provider",
          "resident_provider": "Test Resident Provider"
      },
      "cpt": [
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "SU002",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          }
      ]
    }
  ]
}
`)

func test() {
	// myMap := map[string]interface{}{
	// 	"name":  "John Doe",
	// 	"age":   30,
	// 	"email": "johndoe@example.com",
	// }

	var rules map[string]interface{}
	err := json.Unmarshal(requestData2, &rules)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(rules)
	if err != nil {
		// handle error
	}

	result := sjson.GetBytes(jsonData, "patient.dob")
	if result.Exists() {
		dob := result.Value()
		fmt.Println(dob.(string))
	}

	result = sjson.GetBytes(jsonData, "coding.#.em")
	if result.Exists() {
		em := result.Value()
		fmt.Println(em)
	}

	result = sjson.GetBytes(jsonData, "coding.#.em.code")
	if result.Exists() {
		code := result.Value()
		fmt.Println(code)
	}

	result = sjson.GetBytes(jsonData, "coding.#.cpt.#.code")
	if result.Exists() {
		cpt := result.Value()
		fmt.Println(cpt)
	}
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
		for i := range val.Array() {
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
