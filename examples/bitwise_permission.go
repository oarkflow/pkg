package main

import (
	"fmt"
)

// Define constants for HTTP methods
const (
	GET    = 1 << iota // 1
	POST               // 2
	PUT                // 4
	DELETE             // 8
)

// Define permissions for users in the company
var companyPermissions = map[string]map[string]map[string]map[string]int{
	"companyA": {
		"admin": {
			"IT": {
				"/users":    GET | POST | PUT | DELETE, // Full access
				"/posts":    GET | POST | PUT | DELETE,
				"/comments": GET | POST | PUT | DELETE,
			},
			"HR": {
				"/users":    GET | POST,
				"/posts":    GET | POST,
				"/comments": GET | POST,
			},
		},
		"user": {
			"IT": {
				"/users":    GET | POST,
				"/posts":    GET | POST,
				"/comments": GET | POST,
			},
			"HR": {
				"/users":    GET,
				"/posts":    GET,
				"/comments": GET | POST,
			},
		},
	},
	"companyB": {
		"admin": {
			"IT": {
				"/users":    GET | POST | PUT | DELETE, // Full access
				"/posts":    GET | POST | PUT | DELETE,
				"/comments": GET | POST | PUT | DELETE,
			},
			"HR": {
				"/users":    GET | POST,
				"/posts":    GET | POST,
				"/comments": GET | POST,
			},
		},
		"user": {
			"IT": {
				"/users":    GET | POST,
				"/posts":    GET,
				"/comments": GET | POST,
			},
			"HR": {
				"/users":    GET,
				"/posts":    GET,
				"/comments": GET | POST,
			},
		},
	},
	// Add more companies and their permissions as needed
}

// Function to check if a user in a company has permission for a specific URL and method
func hasCompanyPermission(company, userRole, department, url string, method int) bool {
	companyPermissions, found := companyPermissions[company]
	if !found {
		return false // Company not found
	}
	userPermissions, found := companyPermissions[userRole]
	if !found {
		return false // Role not found
	}
	departmentPermissions, found := userPermissions[department]
	if !found {
		return false // Department not found
	}
	permissions, found := departmentPermissions[url]
	if !found {
		return false // URL not found
	}
	return permissions&method != 0
}

func main() {
	company := "companyA"
	userRole := "user"
	department := "IT"
	url := "/posts"
	method := GET

	if hasCompanyPermission(company, userRole, department, url, method) {
		fmt.Printf("User in %s with role %s in department %s has permission to access %s with method %d\n", company, userRole, department, url, method)
	} else {
		fmt.Printf("User in %s with role %s in department %s does not have permission to access %s with method %d\n", company, userRole, department, url, method)
	}
}

/*
package main

import (
	"fmt"
)

// Define constants for HTTP methods
const (
	GET = 1 << iota // 1
	POST            // 2
	PUT             // 4
	DELETE          // 8
)

// Define permissions for users
var userPermissions = map[string]map[string]map[string]int{
	"companyA": {
		"client1": {
			"service1": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
			"service2": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
		},
		"client2": {
			"service1": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
			"service2": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
		},
	},
	"companyB": {
		"client3": {
			"service1": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
			"service2": {
				"admin":  GET | POST | PUT | DELETE,
				"editor": GET | POST | PUT,
				"viewer": GET,
			},
		},
	},
	// Add more companies, clients, and services as needed
}

// Function to check if a user has permission for a specific operation
func hasPermission(company, client, service, userRole string, method int) bool {
	permissions, found := userPermissions[company][client][service][userRole]
	if !found {
		return false // Permissions not found
	}
	return permissions&method != 0
}

func main() {
	company := "companyA"
	client := "client1"
	service := "service1"
	userRole := "admin"
	url := "/clients/123" // Example URL
	method := GET

	if hasPermission(company, client, service, userRole, method) {
		fmt.Printf("User with role %s in company %s, client %s, and service %s has permission to access %s with method %d\n", userRole, company, client, service, url, method)
	} else {
		fmt.Printf("User with role %s in company %s, client %s, and service %s does not have permission to access %s with method %d\n", userRole, company, client, service, url, method)
	}
}

package main

import (
	"fmt"
	"net/http"
)

// Define constants for HTTP methods
const (
	GET = 1 << iota // 1
	POST            // 2
	PUT             // 4
	DELETE          // 8
)

// Define permissions for users
var userPermissions = map[string]map[string]map[string]map[string]int{
	"companyA": {
		"client1": {
			"service1": {
				"admin":  {"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
				"editor": {"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
				"viewer": {"/users": GET, "/posts": GET},
			},
			"service2": {
				"admin":  {"/comments": GET | POST | PUT | DELETE},
				"editor": {"/comments": GET | POST | PUT},
				"viewer": {"/comments": GET},
			},
		},
		"client2": {
			"service1": {
				"admin":  {"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
				"editor": {"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
				"viewer": {"/users": GET, "/posts": GET},
			},
			"service2": {
				"admin":  {"/comments": GET | POST | PUT | DELETE},
				"editor": {"/comments": GET | POST | PUT},
				"viewer": {"/comments": GET},
			},
		},
	},
	"companyB": {
		"client3": {
			"service1": {
				"admin":  {"/users": GET | POST | PUT | DELETE, "/posts": GET | POST | PUT | DELETE},
				"editor": {"/users": GET | POST | PUT, "/posts": GET | POST | PUT},
				"viewer": {"/users": GET, "/posts": GET},
			},
			"service2": {
				"admin":  {"/comments": GET | POST | PUT | DELETE},
				"editor": {"/comments": GET | POST | PUT},
				"viewer": {"/comments": GET},
			},
		},
	},
	// Add more companies, clients, and services as needed
}

// Middleware function to check permission before serving HTTP request
func checkPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract company, client, service, user role, and URL from the request context
		company := r.Header.Get("X-Company")
		client := r.Header.Get("X-Client")
		service := r.Header.Get("X-Service")
		userRole := r.Header.Get("X-UserRole")
		url := r.URL.Path

		// Check if user has permission for the requested URL and method
		method := getMethodMask(r.Method)
		if !hasPermission(company, client, service, userRole, url, method) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// If permission granted, call the next handler
		next.ServeHTTP(w, r)
	})
}

// Function to map HTTP method to its bitmask
func getMethodMask(method string) int {
	switch method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	default:
		return 0 // Unsupported method
	}
}

// Function to check if a user has permission for a specific operation
func hasPermission(company, client, service, userRole, url string, method int) bool {
	// Check if the company exists
	companyPermissions, found := userPermissions[company]
	if !found {
		return false // Company not found
	}

	// Check if the client exists
	clientPermissions, found := companyPermissions[client]
	if !found {
		return false // Client not found
	}

	// Check if the service exists
	servicePermissions, found := clientPermissions[service]
	if !found {
		return false // Service not found
	}

	// Check if the user role exists
	userRolePermissions, found := servicePermissions[userRole]
	if !found {
		return false // User role not found
	}

	// Check if the URL exists
	permissions, found := userRolePermissions[url]
	if !found {
		return false // URL not found
	}

	// Check if the method is allowed for the URL
	return permissions&method != 0
}

// Example handler function
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Setup routes
	http.HandleFunc("/hello", helloHandler)

	// Use middleware to check permission before serving HTTP requests
	http.Handle("/", checkPermission(http.DefaultServeMux))

	// Start the server
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

*/
