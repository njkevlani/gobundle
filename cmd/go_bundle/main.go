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

type visitedEmelements struct {
	f map[ast.FuncDecl]bool
	s map[ast.StructType]bool
	v map[ast.Expr]bool
}

func main() {
	fileName := "./test_files/test.go"
	res := &ast.File{}
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var mainFunc *ast.FuncDecl
	for _, f := range inFile.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name == "main" {
			mainFunc = fn
		}
	}
	inFile.Decls = append(inFile.Decls, mainFunc)
	var buf bytes.Buffer
	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 8}
	err = printConfig.Fprint(&buf, token.NewFileSet(), inFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf.Bytes()))
	if mainFunc == nil {
		log.Fatal("No main fucntion found in file", fileName)
	}
	visited := &visitedEmelements{
		make(map[ast.FuncDecl]bool),
		make(map[ast.StructType]bool),
		make(map[ast.Expr]bool),
	}
	dfs(mainFunc, res, visited)
}

func dfs(f *ast.FuncDecl, res *ast.File, visited *visitedEmelements) {
	fmt.Printf("f.Name.Name = %+v\n", f.Name.Name)
	for _, i := range f.Body.List {
		fmt.Printf("i = %T\n", i)
	}
}
