package main

import (
	"fmt"

	"github.com/oarkflow/pkg/qr"
)

var (
	paymentData = map[string]map[string]any{
		"esewa": {
			"eSewa_id": "9805832689",
			"name":     "ANITA KUMARI DAS BANIYA",
		},
		/*"bank": {
			"bankCode":      "GLBBNPKA",
			"accountName":   "ANITA KUMARI DAS BANIYA",
			"accountNumber": "334070100044",
			"amount":        "100",
			"remarks":       "This",
		},*/
	}
)

func main() {
	for _, detail := range paymentData {
		q, err := qr.Encode(detail)
		if err != nil {
			panic(err)
		}
		fmt.Println(q.Base64PNG())
		/*err = q.SaveAsPNG(fmt.Sprintf("%s.png", id))
		if err != nil {
			panic(err)
		}*/
	}
}
