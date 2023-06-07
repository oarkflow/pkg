package main

import (
	"fmt"

	"github.com/oarkflow/pkg/filter"
	"github.com/oarkflow/pkg/utils"
)

func multiply(a, b int) int { return a * b }
func triple(a int) int      { return a * 3 }

type FruitRank struct {
	Fruit string
	Rank  uint64
}

func main() {
	fmt.Println(utils.Unique([]string{"abc", "cde", "efg", "efg", "abc", "cde"}))
	fmt.Println(utils.Unique([]int{1, 1, 2, 2, 3, 3, 4}))

	fruits := []FruitRank{
		{
			Fruit: "Strawberry",
			Rank:  1,
		},
		{
			Fruit: "Raspberry",
			Rank:  2,
		},
		{
			Fruit: "Blueberry",
			Rank:  3,
		},
		{
			Fruit: "Blueberry",
			Rank:  3,
		},
		{
			Fruit: "Strawberry",
			Rank:  1,
		},
	}
	fmt.Println(utils.GroupBy(fruits, func(t1 FruitRank) any {
		return t1.Rank
	}))
	fmt.Println(utils.GroupByK(fruits, func(v FruitRank) any {
		return v.Rank
	}, func(t FruitRank) any {
		return t.Rank * 3
	}))
	fmt.Println(utils.Unique(fruits))
}

func Apply() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	factorial := filter.Reduce(a, multiply, 1)
	fmt.Println(factorial)

	filter.ApplyInPlace(a, triple)
	fmt.Println(a)
}
