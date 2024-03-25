package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func main() {
	updated()
}

func updated() {
	et, err := casbin.NewEnforcer("updated-model.conf", "updated-policy.csv")
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}
	fmt.Println(et.Enforce("alice", "companyA", "/code", "write"))               // true
	fmt.Println(et.Enforce("alice", "serviceX", "/check", "read"))               // false
	fmt.Println(et.Enforce("bob", "companyA", "/users", "create-user"))          // true
	fmt.Println(et.Enforce("bob", "companyA", "/code", "write"))                 // true
	fmt.Println(et.Enforce("carol", "companyA", "/companies", "create-company")) // true

}
