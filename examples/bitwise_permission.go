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
