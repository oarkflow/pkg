package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"reflect"
	"strings"
)

type Argument struct {
	Name string
	Type string
}

type Function struct {
	Name      string
	Arguments []string
}

func ma1n() {
	code := `true`
	node, err := parser.ParseExpr(code)
	if err != nil {
		fmt.Println("Error parsing code:", err)
		return
	}
	parser1(node)
}

func parser1(node ast.Node) {
	switch a := node.(type) {
	case *ast.BasicLit:
		fmt.Println("*ast.BasicLit:", a)
	case *ast.Ident:
		fmt.Println("*ast.Ident:", a)
	case *ast.BinaryExpr:
		fmt.Println("*ast.BinaryExpr:", a)
	case *ast.CallExpr:
		fmt.Println("*ast.CallExpr:", a)
	case *ast.BadExpr:
		fmt.Println("*ast.BadExpr:", a)
	case *ast.CompositeLit:
		fmt.Println("*ast.BadExpr:", a)
	case *ast.KeyValueExpr:
		fmt.Println("*ast.KeyValueExpr:", a)
	case *ast.SelectorExpr:
		fmt.Println("End", a.Sel)
		parser1(a.X)
	default:
		fmt.Println("Default:", a, reflect.TypeOf(a))
	}
}

func main() {
	code := "func1(a, true)"
	node, err := parser.ParseExpr(code)
	if err != nil {
		fmt.Println("Error parsing code:", err)
		return
	}
	var functions []Function
	findFunctions(node, &functions)
	fmt.Println(functions)
}

func findFunctions(node ast.Node, functions *[]Function) {
	switch n := node.(type) {
	case *ast.CallExpr:
		*functions = append(*functions, Function{
			Name:      extractFunctionName(n.Fun),
			Arguments: extractArgs(n.Args),
		})
		for _, arg := range n.Args {
			findFunctions(arg, functions)
		}
	}
}

func extractFunctionName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return extractFunctionName(e.Sel)
	}
	return ""
}

func extractArgs(args []ast.Expr) []string {
	var result []string
	for _, arg := range args {
		switch a := arg.(type) {
		case *ast.BasicLit:
			result = append(result, a.Value)
		case *ast.Ident:
			result = append(result, a.Name)
		case *ast.BinaryExpr:
			result = append(result, extractArgs([]ast.Expr{a.X, a.Y})...)
		case *ast.CallExpr:
			result = append(result, extractFunctionCall(a))
		}
	}
	return result
}

func extractFunctionCall(call *ast.CallExpr) string {
	var result []string
	result = append(result, fmt.Sprintf("%s(%s)", extractFunctionName(call.Fun), strings.Join(extractArgs(call.Args), ", ")))
	return strings.Join(result, ", ")
}
