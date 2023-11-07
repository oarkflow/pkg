package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/doc"
)

func main() {
	start := time.Now()
	err := doc.Replace("test.docx", map[string]interface{}{
		"name": "Sujit Baniya",
	})
	if err != nil {
		panic(err)
	}
	/*err := docx.PrepareDocxToFile("test.docx", map[string]interface{}{
		"name": "Sujit Baniya",
	}, "test-filled.docx")
	if err != nil {
		panic(err)
	}*/
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
}
