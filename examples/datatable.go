package main

import (
	"fmt"

	"github.com/oarkflow/pkg/datatable"
)

func main() {
	var rows = []map[string]interface{}{
		{"id": 5, "code": "BJS", "name": "CN", "money": 1.23},
		{"id": 2, "code": "BJS", "name": "CN", "money": 2.21},
		{"id": 3, "code": "SHA", "name": "CN", "money": 1.26},
		{"id": 4, "code": "NYC", "name": "US", "money": 3.99},
		{"id": 7, "code": "MEL", "name": "US", "money": 3.99},
		{"id": 1, "code": "", "name": "CN", "money": 2.99},
	}

	dt := datatable.New(rows)

	table := dt.GroupBy("code,money")
	for i, row := range table.Rows {
		fmt.Println(i, row)
	}
}
