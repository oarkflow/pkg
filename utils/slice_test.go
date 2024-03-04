package utils

import (
	"testing"
)

var (
	s1 = []string{"apple", "banana", "orange"}
	s2 = []string{"banana", "orange", "grape"}
	s3 = []string{"orange", "grape", "kiwi"}
)

func BenchmarkIntersection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Intersection(s1, s2, s3)
	}
}
