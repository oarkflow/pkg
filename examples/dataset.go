package main

import (
	"encoding/json"
	"fmt"

	"github.com/oarkflow/pkg/rule"
)

func main() {
	ruleCheck()
	/*df := dataset.LoadRecords(
		[][]string{
			{"A", "B", "C", "D"},
			{"a", "4", "5.1", "true"},
			{"k", "5", "7.0", "true"},
			{"k", "4", "6.0", "true"},
			{"a", "2", "7.1", "false"},
		},
		dataset.DetectTypes(true),
	)
	fmt.Println(df.GroupBy("A", "D").GetGroups())*/
}

func ruleCheck() {
	data := `[{"code": "A000", "desc": "Cholera due to Vibrio cholerae 01, biovar cholerae"}, {"code": "A001", "desc": "Cholera due to Vibrio cholerae 01, biovar eltor"}, {"code": "A009", "desc": "Cholera, unspecified"}, {"code": "A0100", "desc": "Typhoid fever, unspecified"}, {"code": "A0101", "desc": "Typhoid meningitis"}, {"code": "A0102", "desc": "Typhoid fever with heart involvement"}, {"code": "A0103", "desc": "Typhoid pneumonia"}, {"code": "A0104", "desc": "Typhoid arthritis"}, {"code": "A0105", "desc": "Typhoid osteomyelitis"}]`
	var d rule.Data
	json.Unmarshal([]byte(data), &d)
	r := rule.New()
	r.And(rule.NewCondition("code", rule.IN, []string{"A000", "A001"}))
	fmt.Println(r.Apply(d))
}
