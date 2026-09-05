package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// noZapFatalAnalyzer запрещает zap.Logger.Fatal, zap.Logger.Fatalw,
// zap.Logger.Panic, zap.Logger.Panicw.
var noZapFatalAnalyzer = &analysis.Analyzer{
	Name: "nozapfatal",
	Doc:  "запрещает использование zap.Logger.Fatal, zap.Logger.Fatalw, zap.Logger.Panic, zap.Logger.Panicw",
	Run:  runNoZapFatal,
}

func runNoZapFatal(pass *analysis.Pass) (interface{}, error) {
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

			switch sel.Sel.Name {
			case "Fatal", "Fatalw", "Panic", "Panicw":
				if typ := pass.TypesInfo.TypeOf(sel.X); typ != nil {
					if strings.Contains(typ.String(), "zap.Logger") {
						pass.Reportf(call.Pos(), "zap.Logger.%s запрещён, используйте возврат ошибки", sel.Sel.Name)
					}
				}
			}

			return true
		})
	}
	return nil, nil
}
