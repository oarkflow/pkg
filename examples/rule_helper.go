package main

import (
	"fmt"
)

func searchDeeplyNestedSlice(nestedSlice []interface{}, targetSlice []interface{}) bool {
	targetMap := make(map[interface{}]struct{})
	for _, target := range targetSlice {
		targetMap[target] = struct{}{}
	}

	for _, element := range nestedSlice {
		switch v := element.(type) {
		case []interface{}:
			if searchDeeplyNestedSlice(v, targetSlice) {
				return true
			}
		default:
			if _, found := targetMap[v]; found {
				return true
			}
		}
	}
	return false
}

func main() {
	nestedSlice := []interface{}{
		"apple",
		"banana",
		[]interface{}{"orange", "grape"},
		[]interface{}{
			"kiwi",
			[]interface{}{"pineapple", "mango"},
		},
	}

	targetSlice := []interface{}{"mango", "melon"}

	result := searchDeeplyNestedSlice(nestedSlice, targetSlice)
	fmt.Println(result)
}
