package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/docx"
	"github.com/oarkflow/pkg/evaluate"
	"github.com/oarkflow/pkg/gender"
)

func main() {
	evaluate.AddCustomOperator("as_gender", func(ctx evaluate.EvalContext) (interface{}, error) {
		if ctx.ArgCount() == 2 {
			word, err := ctx.Arg(0)
			if err != nil {
				return nil, err
			}
			gen, err := ctx.Arg(1)
			if err != nil {
				return nil, err
			}
			return gender.Convert(fmt.Sprint(word), fmt.Sprint(gen)), nil
		} else if ctx.ArgCount() == 3 {
			word, err := ctx.Arg(0)
			if err != nil {
				return nil, err
			}
			gen, err := ctx.Arg(1)
			if err != nil {
				return nil, err
			}
			married, err := ctx.BooleanArg(2)
			if err != nil {
				return nil, err
			}
			return gender.Convert(fmt.Sprint(word), fmt.Sprint(gen), married), nil
		}
		return "", nil
	})
	start := time.Now()
	err := docx.PrepareDocxToFile("test.docx", map[string]interface{}{
		"name": "Sujit Baniya",
		"address": map[string]any{
			"city": "Kathmandu",
		},
		"gender": "male",
		"company": map[string]any{
			"name": "Orgware Construct Pvt. Ltd",
		},
		"position":   "Associate Developer",
		"start_date": "2021-09-01",
		"end_date":   "2022-09-30",
	}, "test-filled.docx")
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
}
