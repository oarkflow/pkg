package main

import (
	"fmt"

	"github.com/oarkflow/pkg/dataset"
)

func main() {
	rows := []map[string]interface{}{
		{"id": 5, "code": "BJS", "name": "CN", "money": 1.23},
		{"id": 2, "code": "BJS", "name": "CN", "money": 2.21},
		{"id": 3, "code": "SHA", "name": "CN", "money": 1.26},
		{"id": 4, "code": "NYC", "name": "US", "money": 3.99},
		{"id": 7, "code": "MEL", "name": "US", "money": 3.99},
		{"id": 1, "code": nil, "name": "CN", "money": 2.99},
	}
	rows2 := []map[string]interface{}{
		{"emp_id": 5, "salary": 12000},
	}

	dt := dataset.LoadMaps(rows)
	dt2 := dataset.LoadMaps(rows2)
	fmt.Println(dt.InnerJoin(dt2, dataset.MergeBy{Left: "id", Right: "emp_id"}))
}
