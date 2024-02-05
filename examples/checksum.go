package main

import (
	"fmt"

	"github.com/oarkflow/pkg/checksum"
)

func main() {
	data := []byte(`This is test`)
	key := "r28GYSTF9oJeiXvnkIDLLqu8RGWg3VUS"
	hasher, err := checksum.New64([]byte(key))
	if err != nil {
		panic(err)
	}
	fmt.Println(checksum.MakeSum(hasher, data))
	token := "5468697320697320746573743335e5557ee8c8e4"
	verified := checksum.VerifySum(hasher, data, token)
	fmt.Println(verified)
}
