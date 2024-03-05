package main

import (
	"fmt"

	"github.com/oarkflow/pkg/bitwise"
)

func main() {
	p := bitwise.Factory([]string{"user", "verified", "admin"})

	user := p.Serialize([]string{"user"})
	verified := p.Serialize([]string{"user", "verified"})
	admin := p.Serialize([]string{"user", "admin"})

	fmt.Println(p.Has(user, "user"))
	fmt.Println(p.Has(user, "admin"))
	fmt.Println(p.Has(verified, "verified"))
	fmt.Println(p.Has(verified, "admin"))
	fmt.Println(p.Has(admin, "admin"))

	// add permissions
	fmt.Println(p.Has(user, "verified"))
	user = p.Add(user, "verified")
	fmt.Println(p.Has(user, "verified"))

	// remove permissions
	fmt.Println(p.Has(verified, "verified"))
	verified = p.Remove(verified, "verified")
	fmt.Println(p.Has(verified, "verified"))
}
