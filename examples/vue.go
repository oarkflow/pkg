package main

import (
	"fmt"

	"github.com/oarkflow/pkg/vue"
)

func main() {
	price := vue.Ref(2)
	quantity := vue.Ref(1000)
	revenue := vue.Computed(func() int {
		return price.GetValue() * quantity.GetValue()
	})
	vue.WatchEffect(func() {
		fmt.Println("revenue:", revenue.GetValue())
	})
}
