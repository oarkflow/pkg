package main

import (
	"encoding/json"
	"fmt"

	"github.com/oarkflow/pkg/rule"
	"github.com/oarkflow/pkg/sliceutil"
)

func main() {
	groupBy()
}

type Person struct {
	Name    string
	Age     int
	Country string
}

func groupBy() {
	people := []Person{
		{Name: "Alice", Age: 30, Country: "USA"},
		{Name: "Bob", Age: 25, Country: "Canada"},
		{Name: "Charlie", Age: 30, Country: "USA"},
		{Name: "David", Age: 35, Country: "Canada"},
	}
	persons := []map[string]any{
		{"name": "Alice", "age": 30, "country": "USA"},
		{"name": "Bob", "age": 25, "country": "Canada"},
		{"name": "Charlie", "age": 30, "country": "USA"},
		{"name": "David", "age": 35, "country": "Canada"},
	}
	employees := []map[string]any{
		{"name": "Alice", "country": "USA"},
		{"name": "Bob", "country": "Canada"},
		{"name": "Charlie", "country": "USA"},
		{"name": "David", "country": "Canada"},
	}
	// Use GroupBy with Person type and key function
	fmt.Println(sliceutil.GroupBy(people, "Age", "Country"))
	fmt.Println(sliceutil.GroupBy(persons, "age", "country"))
	fmt.Println(sliceutil.GroupBy(employees, "country"))
}

func ruleCheck() {
	data := `[{"code": "A000", "desc": "Cholera due to Vibrio cholerae 01, biovar cholerae"}, {"code": "A001", "desc": "Cholera due to Vibrio cholerae 01, biovar eltor"}, {"code": "A009", "desc": "Cholera, unspecified"}, {"code": "A0100", "desc": "Typhoid fever, unspecified"}, {"code": "A0101", "desc": "Typhoid meningitis"}, {"code": "A0102", "desc": "Typhoid fever with heart involvement"}, {"code": "A0103", "desc": "Typhoid pneumonia"}, {"code": "A0104", "desc": "Typhoid arthritis"}, {"code": "A0105", "desc": "Typhoid osteomyelitis"}]`
	var d rule.Data
	json.Unmarshal([]byte(data), &d)
	r := rule.New()
	r.And(rule.NewCondition("code", rule.IN, []string{"A000", "A001"}))
	fmt.Println(r.Apply(d))
}
