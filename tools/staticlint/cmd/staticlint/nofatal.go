package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// noFatalAnalyzer запрещает log.Fatal, log.Fatalf, log.Fatalln, zap.Fatal.
var noFatalAnalyzer = &analysis.Analyzer{
	Name: "nofatal",
	Doc:  "запрещает использование log.Fatal, log.Fatalf, log.Fatalln, zap.Fatal",
	Run:  runNoFatal,
}

func runNoFatal(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.Contains(filename, "go-build") || strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkg, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			if pkg.Name == "log" {
				switch sel.Sel.Name {
				case "Fatal", "Fatalf", "Fatalln":
					pass.Reportf(call.Pos(), "log.%s запрещён, используйте возврат ошибки", sel.Sel.Name)
				}
			}

			return true
		})
	}
	return nil, nil
}
