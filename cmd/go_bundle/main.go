package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
)

type visitor struct {
	// TODO Think if I need pointers, will pointer work correctly in my use case?
	res *ast.File
	f   map[*ast.FuncDecl]bool
	pv  map[*ast.GenDecl]bool // package variables
	s   map[*ast.StructType]bool
	v   map[*ast.Expr]bool
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}

	if c, ok := n.(*ast.CallExpr); ok {
		fmt.Println(c)
	}
	s, ok := n.(*ast.SelectorExpr)
	if !ok {
		return v
	}
	x, okX := s.X.(*ast.Ident)
	if !okX {
		return v
	}
	// sel, okSel := s.Sel.(*ast.Ident)
	// if !okSel {
	// 	return v
	// }
	if x.Obj == nil {
		fmt.Printf("need to get func %s in package %s\n", s.Sel, s.X)
	} else {
		fmt.Printf("need to get func %s on struct %s\n", s.Sel, s.X)
	}
	return v
}

func main() {
	fileName := "./test_files/test.go"
	res := &ast.File{}
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var walkFunc *ast.FuncDecl
	for _, f := range inFile.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name == "main" {
			walkFunc = fn
		}
	}
	inFile.Decls = append(inFile.Decls, walkFunc)
	fmt.Println(astToString(inFile))

	if walkFunc == nil {
		log.Fatal("No main function found in file", fileName)
	}
	v := visitor{res: res}
	ast.Walk(v, walkFunc)
}

func astToString(inFile ast.Node) string {
	var buf bytes.Buffer
	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 8}
	err := printConfig.Fprint(&buf, token.NewFileSet(), inFile)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

// TODO Look how goimports add missing imports
