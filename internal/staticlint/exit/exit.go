package exit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer defines the No Exit analyzer.
var Analyzer *analysis.Analyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  "reports calls to os.Exit in the func main of package main.",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case (*ast.FuncDecl):
				if x.Name.Name != "main" {
					return false
				}
			case (*ast.CallExpr):
				doesExit(pass, x)
			case (*ast.Ident):
				isPanic(pass, x)
			}
			return true
		})
	}
	return nil, nil
}

func isPanic(pass *analysis.Pass, ident *ast.Ident) {
	if ident.Name == "panic" {
		pass.Reportf(ident.NamePos, "panic used in main func of package main")
	}
}

func doesExit(pass *analysis.Pass, call *ast.CallExpr) {
	typedCall, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	typedPkg, ok := typedCall.X.(*ast.Ident)
	if !ok {
		return
	}

	switch {
	case typedPkg.Name == "os" && typedCall.Sel.Name == "Exit":
		pass.Reportf(typedCall.Sel.NamePos, "os.Exit used in main func of package main")
	case typedPkg.Name == "log" && strings.HasPrefix(typedCall.Sel.Name, "Fatal"):
		pass.Reportf(typedCall.Sel.NamePos, "log.Fatal or log.Fatalf used in main func of package main")
	}
}
