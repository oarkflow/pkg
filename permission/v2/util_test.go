package v2

import (
	"testing"
)

var patterns = []string{
	"/coding/:start::POST",
	"/encounter::GET",
	"/coding/:wid/:eid/open::GET",
	"/coding/done::POST",
}

var values = []string{
	"/coding/done::POST",
	"/coding/true::POST",
	"/encounter::GET",
	"/coding/123/456/open::GET",
}

func BenchmarkMatchResource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			for _, value := range values {
				MatchResource(value, pattern)
			}
		}
	}
}
