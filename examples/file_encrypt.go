package main

import (
	"github.com/oarkflow/pkg/file"
)

func main() {
	err := file.Decrypt("logo.jpg", []byte("test123"))
	if err != nil {
		panic(err)
	}
}
