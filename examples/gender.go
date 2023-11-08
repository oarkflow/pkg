package main

import (
	"fmt"

	"github.com/oarkflow/pkg/evaluate"
	"github.com/oarkflow/pkg/gender"
)

func main() {
	d := map[string]interface{}{
		"name": "Sujit Baniya",
		"address": map[string]any{
			"city": "Kathmandu",
		},
		"gender": "male",
		"company": map[string]any{
			"name": "Orgware Construct Pvt. Ltd",
			"A":    3,
			"B":    4,
		},
		"position":   "Associate Developer",
		"start_date": "2021-09-01",
		"end_date":   "2022-09-30",
	}
	p, err := evaluate.Parse(`company.A < company.B`, true)
	if err != nil {
		panic(err)
	}
	pt := evaluate.NewEvalParams(d)
	fmt.Println(p.Eval(pt))
	fmt.Println(gender.Convert("his", "male"))
	fmt.Println(gender.Convert("his", "female"))
	fmt.Println(gender.Convert("Mr.", "female"))
}
