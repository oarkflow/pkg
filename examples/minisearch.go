package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oarkflow/pkg/minisearch/pkg/store"
	"github.com/oarkflow/pkg/minisearch/pkg/tokenizer"
)

type ICD struct {
	Code string `json:"code" index:"true"`
	Desc string `json:"desc" index:"true"`
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

func main() {
	data := readFileAsMap("icd10_codes.json")
	db := store.New[any](&store.Config{
		DefaultLanguage: tokenizer.ENGLISH,
		TokenizerConfig: &tokenizer.Config{
			EnableStemming:  true,
			EnableStopWords: true,
		},
		IndexKeys: []string{"code", "desc"},
	})
	p := store.InsertBatchParams[any]{
		Documents: data,
		BatchSize: 100,
	}
	errs := db.InsertBatch(&p)
	if len(errs) > 0 {
		panic(errs)
	}

	s := store.SearchParams{
		Query:      "Cholera",
		Properties: []string{},
		BoolMode:   store.AND,
		Limit:      10,
		Relevance: store.BM25Params{
			K: 1.2,
			B: 0.75,
			D: 0.5,
		},
	}
	rs, err := db.Search(&s)
	if err != nil {
		panic(err)
	}
	fmt.Println(rs.Hits)
}
