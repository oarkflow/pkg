package main

import (
	"fmt"
	"time"

	"github.com/oarkflow/pkg/cache"
)

func main() {
	c := cache.New(cache.WithPersist[string, any](true, "migration"), cache.WithJanitorInterval[string, any](10*time.Second))
	c.Set("a", 1)
	time.Sleep(15 * time.Second)
	fmt.Println(c.Get("a"))
}
