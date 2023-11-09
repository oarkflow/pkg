package main

import (
	"errors"
	"fmt"

	"github.com/oarkflow/expr"

	"github.com/oarkflow/pkg/gender"
	"github.com/oarkflow/pkg/timeutil"
)

func main() {
	expr.AddFunction("age", func(params ...any) (any, error) {
		if len(params) != 1 {
			return nil, errors.New("No data provided")
		}
		left := params[0]
		t, err := timeutil.ParseTime(left)
		if err != nil {
			return nil, err
		}
		return timeutil.CalculateToNow(t), err
	})
	d := map[string]interface{}{
		"customer": map[string]any{
			"dob": "1989-04-10",
		},
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
	p, err := expr.Parse(`age(customer.dob)`)
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Eval(d))
	fmt.Println(gender.Convert("his", "male"))
	fmt.Println(gender.Convert("his", "female"))
	fmt.Println(gender.Convert("Mr.", "female"))
}
