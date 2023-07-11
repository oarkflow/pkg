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
	testStruct()
	testString()
}

type ICD struct {
	Code string `json:"code" index:"code"`
	Desc string `json:"desc" index:"desc"`
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

func testMap() {
	data := readFileAsMap("icd10_codes.json")
	db := search.New[any](&search.Config{
		DefaultLanguage: tokenizer.ENGLISH,
		TokenizerConfig: &tokenizer.Config{
			EnableStemming:  true,
			EnableStopWords: true,
		},
		IndexKeys: search.DocFields(data[0]),
	})
	errs := db.InsertBatch(data, 100)
	if len(errs) > 0 {
		panic(errs)
	}

	s, err := db.Search(&search.Params{
		Extra: map[string]any{
			"code": "A000",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)

	s, err = db.Search(&search.Params{
		Query: "Cholera due to Vibrio cholerae",
		Extra: map[string]any{
			"code": "A001",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
	s, err = db.Search(&search.Params{
		Query: "Cholera",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
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
	if len(errs) > 0 {
		for _, err := range errs {
			panic(err)
		}
	}
	start := time.Now()
	s, err := ftsSearch.Search(&search.Params{
		Extra: map[string]any{
			"code": "Z9981",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Hits)
	fmt.Printf("Time to search %s", time.Since(start))
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
	if len(errs) > 0 {
		for _, err := range errs {
			panic(err)
		}
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
