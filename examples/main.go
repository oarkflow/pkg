package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/server"

	"github.com/oarkflow/pkg/permission"
)

func ma1in() {
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
		{"sujit", "companyA", "/restricted", "GET", ""}, // true
		{"sujit", "companyB", "/restricted", "GET", ""}, // false, expected true
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

func main() {
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
			return []string{"sujit", "companyB", string(ctx.Path()), string(ctx.Method()), string(bt)}
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
