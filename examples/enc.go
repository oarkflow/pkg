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
	passWd := "r28GYSTF9oJeiXvnkIDLLqu8RGWg3VUS"
	encrypted := "c03569b9468afb77e86728ddd7d524f2b4a7996d8de7d9ec7bfc1dd71f6b8e6a24fa229a28fde59508357b9bcfef05513ff2a87ce23ce881c83025ed5c1c5f29"
	fmt.Println(enc.Decrypt(encrypted, passWd))
}
