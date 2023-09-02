package main

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osExit",
	Doc:  "check for os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if funExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkgExpr, ok := funExpr.X.(*ast.Ident); ok {
					if fmt.Sprintf("%s.%s", pkgExpr.Name, funExpr.Sel.Name) == "os.Exit" {
						pass.Reportf(x.Pos(), "os.Exit returns")
					}
				}
			}
		}
	}

	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {

			ast.Inspect(file, func(node ast.Node) bool {
				if y, ok := node.(*ast.FuncDecl); ok {
					if y.Name.String() != "main" {
						return false
					}
				}
				if x, ok := node.(*ast.ExprStmt); ok {
					expr(x)
				}
				return true
			})
		}
	}
	return nil, nil
}
