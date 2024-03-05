package radix

import (
	"testing"
)

// Define constants for HTTP methods
const (
	GET    = 1 << iota // 1
	POST               // 2
	PUT                // 4
	DELETE             // 8
)

var (

	// Convert the nested map structure into a Radix Trie
	userPermissions = map[string]map[string]map[string]map[string]map[string]int{
		"companyA": {
			"client1": {
				"service1": {
					"admin":  map[string]int{"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
					"viewer": map[string]int{"/users": GET, "/posts": GET},
				},
				"service2": {
					"admin":  map[string]int{"/comments": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/comments": GET | POST | PUT},
					"viewer": map[string]int{"/comments": GET},
				},
			},
			"client2": {
				"service1": {
					"admin":  map[string]int{"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
					"viewer": map[string]int{"/users": GET, "/posts": GET},
				},
				"service2": {
					"admin":  map[string]int{"/comments": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/comments": GET | POST | PUT},
					"viewer": map[string]int{"/comments": GET},
				},
			},
		},
		"companyB": {
			"client3": {
				"service1": {
					"admin":  map[string]int{"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
					"viewer": map[string]int{"/users": GET, "/posts": GET},
				},
				"service2": {
					"admin":  map[string]int{"/comments": GET | POST | PUT | DELETE},
					"editor": map[string]int{"/comments": GET | POST | PUT},
					"viewer": map[string]int{"/comments": GET},
				},
			},
		},
	}
)

func BenchmarkRadixInsertPermission(b *testing.B) {
	radixTrie := New()
	for i := 0; i < b.N; i++ {
		// Insert permissions into the Radix Trie
		for company, clients := range userPermissions {
			for client, services := range clients {
				for service, roles := range services {
					for role, urls := range roles {
						for url, bitmask := range urls {
							radixTrie.InsertPermission([]string{company, client, service, role, url}, bitmask)
						}
					}
				}
			}
		}
	}
	b.ReportAllocs()
}

func BenchmarkRadixHasPermission(b *testing.B) {
	data := [][]string{
		{"companyA", "client1", "service1", "admin", "/users"},
		{"companyA", "client1", "service2", "admin", "/comments"},
	}
	radixTrie := New()
	for company, clients := range userPermissions {
		for client, services := range clients {
			for service, roles := range services {
				for role, urls := range roles {
					for url, bitmask := range urls {
						radixTrie.InsertPermission([]string{company, client, service, role, url}, bitmask)
					}
				}
			}
		}
	}
	// Convert the nested map structure into a Radix Trie
	for i := 0; i < b.N; i++ {
		for _, keys := range data {
			radixTrie.HasPermission(keys, GET)
		}
	}
	b.ReportAllocs()
}
