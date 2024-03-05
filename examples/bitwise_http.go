package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oarkflow/pkg/radix"
)

const (
	GET    = 1 << iota // 1
	POST               // 2
	PUT                // 4
	DELETE             // 8
)

var methodMask = map[string]int{
	"GET":    GET,
	"POST":   POST,
	"PUT":    PUT,
	"DELETE": DELETE,
}

var (
	radixTrie *radix.Trie
	skipExt   = []string{".ico", ".css", ".js", ".woff", ".woff2", ".jpg", ".jpeg", ".png", ".gif", ".svg"}
)

func init() {
	radixTrie = radix.New()
	userPermissions := map[string]map[string]map[string]map[string]map[string]int{
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

func checkPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		for _, ext := range skipExt {
			if strings.HasSuffix(url, ext) {
				next.ServeHTTP(w, r)
				return
			}
		}
		keys := []string{"companyA", "client1", "service1", "viewer", "/posts"}
		if !radixTrie.HasPermission(keys, getMethodMask("GET")) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Function to map HTTP method to its bitmask
func getMethodMask(method string) int {
	if bit, ok := methodMask[method]; ok {
		return bit
	}
	return 0
}

func main() {
	// Setup middleware to check permissions
	http.Handle("/", checkPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})))

	// Start the server
	fmt.Println("Listening on http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}
