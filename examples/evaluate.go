package main

import (
	"fmt"

	"github.com/oarkflow/pkg/evaluate"
)

func main() {
	p, err := evaluate.Parse("{{'active'}}", true)
	if err != nil {
		panic(err)
	}
	pr := evaluate.NewEvalParams(nil)
	fmt.Println(p.Eval(pr))
}
