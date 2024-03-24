package main

import (
	"fmt"

	"github.com/oarkflow/pkg/filter"
)

type PermissionRule struct {
	Role     string
	Company  string
	Module   string
	Resource string
	Action   string
	Entities string
}

type UserRolePolicy struct {
	User        string
	Role        string
	Company     string
	Module      string
	Entities    string
	HandleChild string
}

type CompanyPolicy struct {
	Parent string
	Child  string
}

type RolePolicy struct {
	Parent string
	Child  string
}

type Request struct {
	User     string
	Company  string
	Module   string
	Resource string
	Action   string
	Entities string
}

func main() {
	pRules := []PermissionRule{
		{"super-admin", "*", "*", "*", "*", "*"},
		{"account-manager", "company-a", "service-a", "*", "*", "*"},
		{"admin", "company-a", "service-a", "/users", "POST", "1,2,3,4,5,6,7,8,9,10"},
		{"admin", "company-a", "service-a", "/users", "PUT", "1,2,3,4,5,6,7,8,9,10"},
		{"admin", "company-a", "service-a", "/facilities", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"coder", "company-a", "service-a", "/open", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"coder", "company-a", "service-a", "/in-progress", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"qa", "company-a", "service-a", "/qa", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"qa", "company-a", "service-a", "/qa-in-progress", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"de", "company-a", "service-a", "/de", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"de", "company-a", "service-a", "/de-in-progress", "GET", "1,2,3,4,5,6,7,8,9,10"},
		{"suspend-manager", "company-a", "service-a", "/suspend", "GET", "1,2,3,4,5,6,7,8,9,10"},
	}
	uRoles := []UserRolePolicy{
		{"userA", "account - manager", "company - a", "service - a", "*", "true"},
		{"sujit", "super-admin", "*", "*", "*", "true"},
		{"userB", "admin", "company - a", "service - a", "*", "true"},
		{"userC", "coder", "company - a", "service - a", "*", "true"},
		{"userD", "coder", "company - a", "service - a", "2,4,7", "true"},
		{"userE", "qa", "company - a", "service - a", "5,8,9", "true"},
	}

	rPolicy := []RolePolicy{
		{"account - manager", "admin"},
		{"admin", "qa"},
		{"admin", "coder"},
		{"admin", "de"},
		{"admin", "suspend - manager"},
	}
	requests := [][]string{
		{"userA", "company-a", "service-a", "/qa", "GET", "1"},
		{"userA", "company-a", "service-a", "/companies", "GET", "1"},
		{"userB", "company-a", "service-a", "/users", "POST", "1"},
		{"userB", "company-a", "service-a", "/qa", "GET", "1"},
		{"userC", "company-a", "service-a", "/open", "GET", "1"},
	}
	fmt.Println(len(pRules), len(rPolicy))
	for _, request := range requests {
		req := sliceToRequest(request)
		filter.Apply(uRoles, func(rp UserRolePolicy) UserRolePolicy {
			if rp.User == req.User &&
				(rp.Company == req.Company || rp.Company == "*") &&
				(rp.Module == req.Module || rp.Module == "*") {

			}
		})
	}
}

func sliceToRequest(req []string) Request {
	return Request{
		User:     req[0],
		Company:  req[1],
		Module:   req[2],
		Resource: req[3],
		Action:   req[4],
		Entities: req[5],
	}
}
