package main

import (
	"fmt"

	"github.com/oarkflow/pkg/pattern"
)

type User struct {
	ID        string
	FirstName string
}

func main() {
	a := 3
	b := 115
	result, err := pattern.
		Match(a, b).
		Case(func(args ...any) (any, error) {
			return 5, nil
		}, 3, 15).
		Case(func(args ...any) (any, error) {
			return 4, nil
		}, pattern.EXISTS, pattern.ANY).
		Default(func(args ...any) (any, error) {
			return 2, nil
		}).
		Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result, err)
}
