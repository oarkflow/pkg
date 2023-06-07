package main

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/server"
	
	"github.com/oarkflow/pkg/permission"
)

func main() {
	perm, err := permission.Default(permission.Config{
		Model:  "model.conf",
		Policy: "policy.csv",
		ParamExtractor: func(c context.Context, ctx *frame.Context) []string {
			bt, _ := json.Marshal(map[string]any{
				"service": "medical-coding",
				"entity":  "work-item",
			})
			return []string{"sujit", "domainHere", "test", "write", string(bt)}
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
	srv.Spin()
}
