package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/server"

	"github.com/oarkflow/pkg/permission"
)

func main() {
	updated()
}

func updated() {
	et, err := permission.Default(permission.Config{
		Model:  "updated-model.conf",
		Policy: "updated-policy.csv",
		ParamExtractor: func(c context.Context, ctx *frame.Context) []string {
			bt, _ := json.Marshal(map[string]any{
				"service":   "medical-coding",
				"entity":    "work-item",
				"entity_id": "1",
			})
			// "user", "company", "url/feature", "method/action", "json attributes"
			return []string{"sujit", "arnet", string(ctx.Path()), string(ctx.Method()), string(bt)}
		},
	})
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}
	fmt.Println(et.GetUserPermissions("userB", "company-a"))
	fmt.Println(et.GetDomainsForUser("userB"))
	fmt.Println(et.GetRolesInDomain("company-a"))
	fmt.Println(et.GetAllUsersInDomainWithRole("company-a"))
	fmt.Println(et.GetUserRole("company-a", "userB"))
	fmt.Println(et.GetModuleRelatedByRole("company-a", "coder"))
	fmt.Println(et.GetImplicitRolesForUser("userB", "company-a"))
	fmt.Println(et.GetImplicitResourcesForUser("userB", "company-a"))
	fmt.Println(et.GetAllNamedActions("p"))
	/*slice := [][]any{
		{"userA", "company-a", "service-a", "/qa", "GET", "1"},        // expected true, actual true
		{"userA", "company-a", "service-a", "/companies", "GET", "1"}, // expected true, actual true
		{"userB", "company-a", "service-a", "/users", "POST", "1"},    // expected true, actual true
		{"userB", "company-a", "service-a", "/qa", "GET", "1"},        // expected true, actual false
		{"userC", "company-a", "service-a", "/open", "GET", "1"},      // expected true, actual true
	}
	for _, rVals := range slice {
		ok, err := et.Enforce(rVals...)
		if err != nil {
			fmt.Println("error: ")
			panic(err)
		}
		fmt.Println("Route Valid", ok)
	}*/
}

func attributes() {
	et, err := permission.Default(permission.Config{
		Model:  "model.conf",
		Policy: "policy.csv",
		ParamExtractor: func(c context.Context, ctx *frame.Context) []string {
			bt, _ := json.Marshal(map[string]any{
				"service":   "medical-coding",
				"entity":    "work-item",
				"entity_id": "1",
			})
			// "user", "company", "url/feature", "method/action", "json attributes"
			return []string{"sujit", "arnet", string(ctx.Path()), string(ctx.Method()), string(bt)}
		},
	})
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}
	slice := [][]any{
		{"sujit", "companyA", "/users", "GET", ""},      // true
		{"sujit", "companyB", "/restricted", "GET", ""}, // false, expected true
		{"sujit", "companyC", "/post", "GET", ""},       // false, expected true
	}
	for _, rVals := range slice {
		ok, err := et.Enforce(rVals...)
		if err != nil {
			fmt.Println("error: ")
			panic(err)
		}
		fmt.Println("Route Valid", ok)
	}
}

func main1() {
	perm, err := permission.Default(permission.Config{
		Model:      "model.conf",
		Policy:     "policy.csv",
		SkipExcept: []string{"GET /"},
		ParamExtractor: func(c context.Context, ctx *frame.Context) []string {
			bt, _ := json.Marshal(map[string]any{
				"service":   "medical-coding",
				"entity":    "work-item",
				"entity_id": "1",
			})
			// "user", "company", "url/feature", "method/action", "json attributes"
			return []string{"sujit", "companyB", string(ctx.Path()), string(ctx.Method()), string(bt)}
		},
	})
	if err != nil {
		fmt.Println("error: ")
		panic(err)
	}
	srv := server.Default(server.WithExitWaitTime(1 * time.Second))
	srv.Use(perm.RoutePermission)
	srv.GET("/", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, "Welcome home")
	})
	srv.GET("/restrict/:id", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, ctx.FullPath())
	})
	srv.POST("/restrict", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, "Access Done")
	})
	srv.NoRoute(func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, "No route")
	})
	srv.Spin()
}
