package main

import (
	"fmt"

	"github.com/oarkflow/pkg/timeutil"
)

func main() {
	datetimeStr := "2079/10/14"
	format := "%Y/%m/%d"

	npTime, err := timeutil.ParseNP(datetimeStr, format)
	if err != nil {
		panic(npTime)
	}
	fmt.Println(npTime)
}
