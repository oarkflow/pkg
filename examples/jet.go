package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/domainutil"
	"github.com/oarkflow/pkg/jet"
)

var (
	data3 = map[string]any{
		"first_name": "Sujit",
		"address": map[string]any{
			"city": "Kathmandu",
		},
	}
	data4 = map[string]any{
		"first_name": "Anita",
		"address": map[string]any{
			"city": "Delhi",
		},
	}
)

func main() {
	fmt.Println(domainutil.Domain("keep.google.co.uk"))
	fmt.Println(domainutil.Subdomain("google.co.uk"))
	fmt.Println(domainutil.SplitDomain("keep.google.co.uk"))
	fmt.Println(domainutil.DomainSuffix("keep.google.co.uk"))
	fmt.Println(domainutil.DomainPrefix("keep.google.co.uk"))
	// country()
}
func jetTest() {
	jetParse()
	jetTemplateParse()
}

func jetParse() {
	start := time.Now()
	fmt.Println(jet.Parse("Hi Mr. {{ address.city }}", data3))
	fmt.Println(jet.Parse("Hi Mr. {{ address.city }}", data4))
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
}

func jetTemplateParse() {
	start := time.Now()
	tmpl, err := jet.NewTemplate("Hi Mr. {{ address.city }}")
	if err != nil {
		panic(err)
	}
	fmt.Println(tmpl.Parse(data3))
	fmt.Println(tmpl.Parse(data4))
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
}
