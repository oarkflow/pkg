package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/jet"
	"github.com/oarkflow/pkg/str"
)

func main() {
	// strFormat()
	jetFormat()
}

func jetFormat() {
	start := time.Now()
	fmt.Println(jet.Sprintf("client:<user.id>:<request_id>", map[string]any{
		"user": map[string]any{
			"id": 2,
		},
		"request_id": 5,
		"file_id":    "1",
	}, &jet.Delims{
		Left:  "<",
		Right: ">",
	}))
	fmt.Println(time.Since(start))
}

func strFormat() {
	start := time.Now()
	str.Printfln("client:%<user.id>.f:<request_id>", map[string]any{
		"user": map[string]any{
			"id": 2,
		},
		"request_id": 5,
		"file_id":    "1",
	})
	str.Printfln("client:<user_id>:<request_id>", map[string]any{
		"user_id":    2,
		"request_id": 5,
		"file_id":    "1",
	})
	fmt.Println(time.Since(start))
}
