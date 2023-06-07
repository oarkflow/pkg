package main

import (
	"fmt"

	"github.com/oarkflow/pkg/dipper"
)

func main() {
	// structData()
	mapData()
}

type Em struct {
	Code             string `json:"code"`
	BillingProvider  string `json:"billing_provider"`
	ResidentProvider string `json:"resident_provider"`
	EncounterUid     int    `json:"encounter_uid"`
	WorkItemUid      int    `json:"work_item_uid"`
}

type Cpt struct {
	Code             string `json:"code"`
	BillingProvider  string `json:"billing_provider"`
	ResidentProvider string `json:"resident_provider"`
	EncounterUid     int    `json:"encounter_uid"`
	WorkItemUid      int    `json:"work_item_uid"`
}

type Request struct {
	Em  Em    `json:"em"`
	Cpt []Cpt `json:"cpt"`
}

func structData() {
	data := Request{
		Em: Em{
			Code:             "001",
			EncounterUid:     1,
			WorkItemUid:      2,
			BillingProvider:  "Test provider",
			ResidentProvider: "Test Resident Provider",
		},
		Cpt: []Cpt{
			{
				Code:             "001",
				EncounterUid:     1,
				WorkItemUid:      2,
				BillingProvider:  "Test provider",
				ResidentProvider: "Test Resident Provider",
			},
			{
				Code:             "OBS01",
				EncounterUid:     1,
				WorkItemUid:      2,
				BillingProvider:  "Test provider",
				ResidentProvider: "Test Resident Provider",
			},
			{
				Code:             "SU002",
				EncounterUid:     1,
				WorkItemUid:      2,
				BillingProvider:  "Test provider",
				ResidentProvider: "Test Resident Provider",
			},
		},
	}
	fmt.Println(dipper.Get(data, "Cpt.[].Code"))
}

func mapData() {
	data := []map[string]any{
		{
			"code":              "001",
			"billing_provider":  "Test provider",
			"resident_provider": "Test Resident Provider",
		},
		{
			"code":              "OBS01",
			"billing_provider":  "Test provider",
			"resident_provider": "Test Resident Provider",
		},
		{
			"code":              "SU002",
			"billing_provider":  "Test provider",
			"resident_provider": "Test Resident Provider",
		},
	}
	fmt.Println(dipper.FilterSlice(data, ".[].code", []any{"OBS01"}))
}
