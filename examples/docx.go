package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/oarkflow/expr"

	"github.com/oarkflow/pkg/docx"
	"github.com/oarkflow/pkg/gender"
	"github.com/oarkflow/pkg/timeutil"
)

func tableParser() {
	replaceMap := [][]docx.TablePlaceholder{
		{
			{
				Key:   "Name",
				Value: "John Doe",
			},
			{
				Key:   "Value",
				Value: "42",
			},
		},
		{
			{
				Key:   "Name",
				Value: "Jane Doe",
			},
			{
				Key:   "Value",
				Value: "43",
			},
		},
	}

	doc, err := docx.Open("/Users/sujit/Sites/pkg/examples/table_template.docx")
	if err != nil {
		panic(err)
		return
	}

	tableReplacer := docx.NewTableReplacer(doc.GetFile(docx.DocumentXml))
	err = tableReplacer.Replace("t1", replaceMap)
	if err != nil {
		panic(err)
		return
	}

	doc.SetFile(docx.DocumentXml, tableReplacer.Bytes())

	err = doc.WriteToFile("./test/out.docx")
	if err != nil {
		panic(err)
		return
	}
}

func main() {
	expr.AddFunction("current_date", func(params ...any) (any, error) {
		return time.Now().Format(time.DateOnly), nil
	})
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
	expr.AddFunction("as_gender", func(params ...any) (any, error) {
		if len(params) == 2 {
			word := params[0]
			gen := params[1]
			if word == nil {
				word = ""
			}
			if gen == nil {
				gen = ""
			}
			return gender.Convert(fmt.Sprint(word), fmt.Sprint(gen)), nil
		} else if len(params) == 3 {
			word := params[0]
			gen := params[1]
			married := params[2]
			if word == nil {
				word = ""
			}
			if gen == nil {
				gen = ""
			}
			if married == nil {
				married = false
			}
			return gender.Convert(fmt.Sprint(word), fmt.Sprint(gen), married.(bool)), nil
		}
		return "", nil
	})
	doc := "/home/sujit/Projects/paramarsha/frontend/public/test.docx"
	start := time.Now()
	placeholders, err := docx.Placeholders(doc)
	if err != nil {
		panic(err)
	}
	fmt.Println(placeholders)
	err = docx.PrepareDocxToFile(doc, map[string]interface{}{
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
		},
		"position":   "Associate Developer",
		"start_date": "2021-09-01",
		"end_date":   "2022-09-30",
	}, "test-filled.docx")
	if err != nil {
		panic(err)
	}
	// ([a-zA-Z_]\w*)\(([^()]|(?R))*\)
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
}
