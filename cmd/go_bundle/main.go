package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"

	"os"
	"path"

	"golang.org/x/tools/imports"
)

type funcDeclCollectorVisitor struct {
	curPkg      string
	funcDeclMap map[string]*ast.FuncDecl
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			editedFuncName := funcDecl.Name.Name
			if fv.curPkg != "" {
				editedFuncName = fmt.Sprintf("%s_%s", fv.curPkg, funcDecl.Name.Name)
			}

			if _, alreadyExists := fv.funcDeclMap[editedFuncName]; alreadyExists {
				log.Fatalf("Function already exists in map. editedFuncName=%s\n", editedFuncName)
			}

			funcDecl.Name.Name = editedFuncName
			fv.funcDeclMap[editedFuncName] = funcDecl
		}
	}

	return fv
}

type visitor struct {
	// TODO Think if I need pointers, will pointer work correctly in my use case?
	res *ast.File
	f   map[string]bool
	pv  map[*ast.GenDecl]bool // package variables
	s   map[*ast.StructType]bool
	v   map[*ast.Expr]bool
	fv  *funcDeclCollectorVisitor
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			// spew.Dump(callExpr)
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				// TODO: fix it. It will work only for "algo" package.
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && pkgIdent.Name == "algo" {
					editedFuncName := pkgIdent.Name + "_" + selectorExpr.Sel.Name
					callExpr.Fun = ast.NewIdent(editedFuncName)

					if !v.f[editedFuncName] {
						v.res.Decls = append(v.res.Decls, v.fv.funcDeclMap[editedFuncName])
						v.f[editedFuncName] = true
					}
				}
			} else if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				funcName := ident.Name
				if !v.f[funcName] {
					v.res.Decls = append(v.res.Decls, v.fv.funcDeclMap[funcName])
					v.f[funcName] = true
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

	pkgs, err := parser.ParseDir(token.NewFileSet(), "./test_files/algo", nil, parser.ParseComments)

	fv := funcDeclCollectorVisitor{funcDeclMap: make(map[string]*ast.FuncDecl)}
	if err != nil {
		log.Fatal(err)
	}
	for pkgName, pkg := range pkgs {
		for _, file := range pkg.Files {
			fv.curPkg = pkgName
			ast.Walk(fv, file)
		}
	}

	fileName := path.Clean(os.Args[1])
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fv.curPkg = ""
	ast.Walk(fv, inFile)

	var mainFunc *ast.FuncDecl
	for _, f := range inFile.Decls {
		if fn, ok := f.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
			mainFunc = fn
			break
		}
	}

	if mainFunc == nil {
		log.Fatal("No main function found in file: ", fileName)
	}

	res := &ast.File{Name: ast.NewIdent("main")}
	res.Decls = append(res.Decls, mainFunc)
	v := visitor{res: res, fv: &fv, f: make(map[string]bool)}
	ast.Walk(v, res)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), res)
	if err != nil {
		log.Fatal(err)
	}

	bytes := buf.Bytes()
	goimportedBytes, err := imports.Process("", bytes, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(goimportedBytes))
}
