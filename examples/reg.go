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
	Arguments []Argument
}

func man() {
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
	code := "func1(arg1, func2(arg2, func3(arg3, a.b)), func4(arg4))"
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

func getSelectExpr(arg ast.Expr, b string) string {
	switch a := arg.(type) {
	case *ast.BasicLit:
		return a.Value + "." + b
	case *ast.Ident:
		return a.Name + "." + b
	case *ast.SelectorExpr:
		return getSelectExpr(a.X, a.Sel.Name+"."+b)
	}
	return b
}

func extractArg(arg ast.Expr) Argument {
	varg := Argument{}
	switch a := arg.(type) {
	case *ast.BasicLit:
		varg.Type = "constant"
		varg.Name = a.Value
	case *ast.Ident:
		if a.Name == "true" || a.Name == "false" {
			varg.Type = "constant"
			varg.Name = a.Name
		} else {
			varg.Type = "variable"
			varg.Name = a.Name
		}
	case *ast.SelectorExpr:
		varg.Type = "variable"
		varg.Name = getSelectExpr(a.X, a.Sel.Name)
	case *ast.CallExpr:
		varg.Type = "function"
		varg.Name = extractFunctionCall(a)
	}
	return varg
}

func extractArgs(args []ast.Expr) []Argument {
	var result []Argument
	for _, arg := range args {
		result = append(result, extractArg(arg))
	}
	return result
}

func extractFunctionCall(call *ast.CallExpr) string {
	var result []string
	var args []string
	for _, a := range extractArgs(call.Args) {
		args = append(args, a.Name)
	}
	result = append(result, fmt.Sprintf("%s(%s)", extractFunctionName(call.Fun), strings.Join(args, ", ")))
	return strings.Join(result, ", ")
}
