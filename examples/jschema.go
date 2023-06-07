package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/oarkflow/pkg/jsonschema"
)

var data = []byte(`
{
    "em": {
        "code": "001",
        "encounter_uid": 1,
        "work_item_uid": 2, 
        "billing_provider": "Test provider",
        "resident_provider": "Test Resident Provider"
    },
    "cpt": [
        {
            "code": "001",
            "billing_provider": "Test provider",
            "resident_provider": "Test Resident Provider"
        },
        {
            "code": "OBS01",
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
`)

func main() {
	ctx := context.Background()
	schemaData, _ := os.ReadFile("schema.json")

	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	errs, err := rs.ValidateBytes(ctx, data)
	if err != nil {
		panic(err)
	}

	if len(errs) > 0 {
		fmt.Println(errs[0].Error())
	}
}
