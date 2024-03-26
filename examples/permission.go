package main

import (
	"fmt"
	"strings"
)

func main() {
	routePatterns := []string{
		"/coding/:start::POST",
		"/encounter::GET",
		"/coding/:wid/:eid/open::GET",
		"/coding/done::POST",
	}

	values := []string{
		"/coding/done::POST",
		"/coding/true::POST",
		"/encounter::GET",
		"/coding/123/456/open::GET",
	}

	for _, value := range values {
		fmt.Printf("Value: %s\n", value)
		matched := false
		for _, pattern := range routePatterns {
			if MatchResource3(value, pattern) {
				matched = true
				fmt.Printf("Matched pattern: %s\n", pattern)
				break
			}
		}
		if !matched {
			fmt.Println("No matching pattern found")
		}
		fmt.Println("----------------------")
	}
}

func MatchResource3(value, pattern string) bool {
	valueParts := strings.Split(value, "/")
	patternParts := strings.Split(pattern, "/")

	if len(valueParts) != len(patternParts) {
		return false
	}

	for i := 0; i < len(valueParts); i++ {
		if patternParts[i] == "*" || patternParts[i] == valueParts[i] {
			continue
		}

		if strings.HasPrefix(patternParts[i], ":") {
			continue
		}

		return false
	}

	return true
}
