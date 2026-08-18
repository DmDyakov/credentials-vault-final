package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// noOsExitAnalyzer запрещает os.Exit вне функции main пакета main.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "запрещает os.Exit вне функции main пакета main",
	Run:  runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.Contains(filename, "go-build") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			isMainFunc := fn.Name.Name == "main" && fn.Recv == nil && pass.Pkg.Name() == "main"

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkg, ok := sel.X.(*ast.Ident)
				if !ok || pkg.Name != "os" || sel.Sel.Name != "Exit" {
					return true
				}

				if !isMainFunc {
					pass.Reportf(call.Pos(), "os.Exit разрешён только в main функции пакета main")
				}

				return true
			})
			return true
		})
	}
	return nil, nil
}
