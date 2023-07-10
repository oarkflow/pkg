package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/oarkflow/pkg/search"
	"github.com/oarkflow/pkg/search/tokenizer"
)

func main() {
	testMap()
	// testStruct()
	// testString()
}

func testStruct() {
	data := readData()
	ftsSearch := search.New[ICD](&search.Config{
		TokenizerConfig: &tokenizer.Config{
			EnableStopWords: true,
			EnableStemming:  true,
		},
	})
	errs := ftsSearch.InsertBatch(data, 2)
	if errs != nil {
		panic(errs)
	}
	start := time.Now()
	s, err := ftsSearch.Search(&search.Params{
		Exact:      true,
		BoolMode:   search.AND,
		Properties: map[string]bool{"Desc": true},
		Offset:     0,
		Limit:      10,
		Extra: map[string]any{
			"code": "O24911",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
	fmt.Printf("Time to search %s", time.Since(start))
}

func testMap() {
	data := readFileAsMap("icd10_codes.json")
	ftsSearch := search.New[any](&search.Config{
		TokenizerConfig: &tokenizer.Config{
			EnableStopWords: true,
			EnableStemming:  true,
		},
	})
	/*for _, d := range data {
		_, err := ftsSearch.Insert(d, tokenizer.ENGLISH)
		if err != nil {
			panic(err)
		}
	}*/
	errs := ftsSearch.InsertBatch(data, 5)
	if errs != nil {
		panic(errs)
	}

	s, err := ftsSearch.Search(&search.Params{
		BoolMode: search.AND,
		Exact:    false,
		Extra: map[string]any{
			"code": "A000",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)

	s, err = ftsSearch.Search(&search.Params{
		Query:    "Cholera due to Vibrio cholerae",
		BoolMode: search.AND,
		Exact:    false,
		Extra: map[string]any{
			"code": "A001",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)

	s, err = ftsSearch.Search(&search.Params{
		Query:    "Cholera",
		BoolMode: search.AND,
		Exact:    false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
}

func testString() {
	data := readFromInt()
	ftsSearch := search.New[int](&search.Config{
		TokenizerConfig: &tokenizer.Config{
			EnableStopWords: true,
			EnableStemming:  true,
		},
	})
	errs := ftsSearch.InsertBatch(data, 2)
	if errs != nil {
		panic(errs)
	}
	start := time.Now()
	s, err := ftsSearch.Search(&search.Params{
		Query: "10",
		Exact: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
	fmt.Printf("Time to search %s", time.Since(start))
}

type ICD struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func readData() (icds []ICD) {
	jsonData, err := os.ReadFile("icd10_codes.json")
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}

	if err := json.Unmarshal(jsonData, &icds); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	return
}

func readFileAsMap(file string) (icds []any) {
	jsonData, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to read json file, error: %v", err)
		return
	}

	if err := json.Unmarshal(jsonData, &icds); err != nil {
		fmt.Printf("failed to unmarshal json file, error: %v", err)
		return
	}
	return
}

func readFromString() []string {
	return []string{
		"Salmonella pneumonia",
		"Diabetes uncontrolled",
	}
}

func readFromInt() []int {
	return []int{
		10,
		100,
		20,
	}
}
