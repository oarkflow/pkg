package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/server"

	"github.com/oarkflow/pkg/permission"
)

func main() {
	e, err := casbin.NewEnforcer("model.conf", "working_policy.csv")
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}
	if len(permission.CasFunc) > 0 {
		for key, fn := range permission.CasFunc {
			e.AddFunction(key, fn)
		}
	}
	slice := [][]any{
		{"alice", "domain1", "data1", "read", `{entity:"work-item"}`},
		{"alice", "domain1", "data2", "read", ""},
		{"alice", "domain1", "data3", "read", ""},
		{"bob", "domain2", "data1", "read", ""},
		{"bob", "domain2", "data2", "read", ""},
		{"bob", "domain2", "data3", "read", ""},
		{"carol", "domain3", "data1", "read", ""},
		{"carol", "domain3", "data2", "read", ""},
		{"carol", "domain3", "data3", "read", ""},
	}
	for _, rVals := range slice {
		ok, err := e.Enforce(rVals...)
		if err != nil {
			fmt.Println("error: ")
			panic(err)
		}
		fmt.Println("Valid", ok)
	}
	et, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}
	if len(permission.CasFunc) > 0 {
		for key, fn := range permission.CasFunc {
			et.AddFunction(key, fn)
		}
	}
	slice = [][]any{
		{"sujit", "companyA", "/restricted", "GET", ""}, // true
		{"sujit", "companyB", "/restricted", "GET", ""}, // true
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

func routeMiddlePermission() {
	perm, err := permission.Default(permission.Config{
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
		fmt.Println("error: ")
		panic(err)
	}
	srv := server.Default()
	srv.Use(perm.RoutePermission)
	srv.GET("/restrict", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, "Access Done")
	})
	srv.POST("/restrict", func(c context.Context, ctx *frame.Context) {
		ctx.JSON(200, "Access Done")
	})
	srv.Spin()
}
