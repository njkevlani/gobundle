package builtinfuncdetector

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/packages"
)

var isBuiltinFunc = make(map[string]bool)

type builtinFuncCollector struct {
	collected map[string]bool
}

func (visitor builtinFuncCollector) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			visitor.collected[funcDecl.Name.Name] = true
		}
	}

	return visitor
}

func init() {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedSyntax}, "builtin")

	if err != nil {
		log.Fatal(err)
	}

	visitor := builtinFuncCollector{isBuiltinFunc}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Walk(visitor, file)
		}
	}
}

// isBuiltinFunc returns true if given string is a function in https://pkg.go.dev/builtin, else false.
func IsBuiltinFunc(funcName string) bool {
	return isBuiltinFunc[funcName]
}
