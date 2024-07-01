package main

import (
	"encoding/binary"
	"fmt"

	"github.com/oarkflow/xid"

	"github.com/oarkflow/pkg/enc"
)

func main() {
	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, uint64(xid.New().Int64()))
	fmt.Println(xid.New().String())
	key := "myverystrongpasswordo32bitlength"
	originalText := "myverystrongpasswordo32bitlength"
	fmt.Println(enc.Encrypt(originalText, key))
}
