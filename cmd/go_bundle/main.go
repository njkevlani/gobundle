package main

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path"
	// "golang.org/x/tools/go/ast/astutil"
)

// func getFuncNameWithPkg() {
// 				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && pkgIdent.Name == "algo" {
// 					editedFuName := ast.NewIdent(pkgIdent.Name + "_" + selectorExpr.Sel.Name)
// 					callExpr.Fun = editedFuName
// 					funcDefinition := v.getFuncDecl()
// 					v.res.Decls = append(v.res.Decls, funcDefinition)
// 				}
// }

type funcDeclCollectorVisitor struct {
	funcDeclMap map[string]*ast.FuncDecl
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			spew.Dump(funcDecl)
		}
	}

	return fv
}

type visitor struct {
	// TODO Think if I need pointers, will pointer work correctly in my use case?
	res *ast.File
	f   map[*ast.FuncDecl]bool
	pv  map[*ast.GenDecl]bool // package variables
	s   map[*ast.StructType]bool
	v   map[*ast.Expr]bool
	fv  *funcDeclCollectorVisitor
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			spew.Dump(callExpr)
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				// TODO: fix it. It will work only for "algo" package.
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && pkgIdent.Name == "algo" {
					editedFuName := ast.NewIdent(pkgIdent.Name + "_" + selectorExpr.Sel.Name)
					callExpr.Fun = editedFuName
					funcDefinition := v.getFuncDecl()
					v.res.Decls = append(v.res.Decls, funcDefinition)
				}
			}
		}
	}

	return v
}

func (v visitor) getFuncDecl() *ast.FuncDecl {
	// TODO: Get *ast.Package and do ast.Walk(pkg, v.fv)
	return nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Arguments required")
	}

	fileName := path.Clean(os.Args[1])
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// spew.Dump(inFile)

	// importList := astutil.Imports(token.NewFileSet(), inFile)
	// spew.Dump(importList)

	// fmt.Println("decls")
	// spew.Dump(inFile.Decls)

	// var walkFunc *ast.FuncDecl

	// for _, decl := range inFile.Decls {
	// 	if fn, ok := decl.(*ast.FuncDecl); ok {
	// 		if fn.Name.Name == "main" {
	// 			fn.Name.Name = "lorem"
	// 			// spew.Dump(fn.Body.List)
	// 		}
	// 		// TODO: Get callExpr out of this and mainupulate it.
	// 		for _, statements := range fn.Body.List {
	// 			if expStatement, ok := statements.(*ast.ExprStmt); ok {
	// 				if callExpr, ok := expStatement.X.(*ast.CallExpr); ok {
	// 					fmt.Println("CallExpr:")
	// 					spew.Dump(callExpr)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// inFile.Decls = append(inFile.Decls, walkFunc)
	// fmt.Println(astToString(inFile))

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

	if mainFunc == nil {
		log.Fatal("No main function found in file", fileName)
	}

	// res := &ast.File{}
	// res.Decls = append(res.Decls, mainFunc)
	v := visitor{res: inFile}
	ast.Walk(v, inFile)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), inFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
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
