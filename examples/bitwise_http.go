package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oarkflow/pkg/radix"
)

// Define constants for HTTP methods
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

// Middleware function to check permission before serving HTTP request
func checkPermission(next http.Handler, rt *radix.Trie) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract keys and method from the request context
		company := r.Header.Get("X-Company")
		client := r.Header.Get("X-Client")
		service := r.Header.Get("X-Service")
		role := r.Header.Get("X-Role")
		url := r.URL.Path
		method := getMethodMask(r.Method)

		// List of static file extensions to skip permission checks for
		staticFileExtensions := []string{".ico", ".css", ".js", ".woff", ".woff2", ".jpg", ".jpeg", ".png", ".gif", ".svg"}

		// Check if the request targets a static file
		for _, ext := range staticFileExtensions {
			if strings.HasSuffix(url, ext) {
				// Skip permission check for static files
				next.ServeHTTP(w, r)
				return
			}
		}

		// Check if user has permission for the requested URL and method
		keys := []string{company, client, service, role, url}
		if !rt.HasPermission(keys, method) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// If permission granted, call the next handler
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
	// Initialize a new Radix Trie
	radixTrie := radix.New()

	// Convert the nested map structure into a Radix Trie
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

	// Setup middleware to check permissions
	http.Handle("/", checkPermission(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}), radixTrie))

	// Start the server
	fmt.Println("Listening on http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}
