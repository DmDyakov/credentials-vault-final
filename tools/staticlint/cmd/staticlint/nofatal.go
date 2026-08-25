package main

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// noFatalAnalyzer запрещает log.Fatal, log.Fatalf, log.Fatalln.
var noFatalAnalyzer = &analysis.Analyzer{
	Name: "nofatal",
	Doc:  "запрещает использование log.Fatal, log.Fatalf, log.Fatalln",
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

			fn, ok := pass.TypesInfo.Uses[sel.Sel].(*types.Func)
			if !ok {
				return true
			}

			pkg := fn.Pkg()
			if pkg == nil || pkg.Path() != "log" {
				return true
			}

			switch sel.Sel.Name {
			case "Fatal", "Fatalf", "Fatalln":
				pass.Reportf(call.Pos(), "log.%s запрещён, используйте возврат ошибки", sel.Sel.Name)
			}

			return true
		})
	}
	return nil, nil
}
