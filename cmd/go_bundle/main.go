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

func getEditedFuncName(pkgName, funcName string) string {
	if pkgName == "" {
		return funcName
	}

	return fmt.Sprintf("%s_%s", pkgName, funcName)
}

type funcDeclCollectorVisitor struct {
	curPkg      string
	funcDeclMap map[string]*ast.FuncDecl
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			editedFuncName := getEditedFuncName(fv.curPkg, funcDecl.Name.Name)

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
	fv       *funcDeclCollectorVisitor
	result   *ast.File
	curPkg   string
	doneFunc map[string]bool
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				// TODO: fix it. It will work only for "algo" package.
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && pkgIdent.Name == "algo" {
					editedFuncName := pkgIdent.Name + "_" + selectorExpr.Sel.Name
					callExpr.Fun = ast.NewIdent(editedFuncName)

					if !v.doneFunc[editedFuncName] {
						funcDecl := v.fv.funcDeclMap[editedFuncName]

						// Add this function in result.
						v.result.Decls = append(v.result.Decls, funcDecl)
						v.doneFunc[editedFuncName] = true
						curPkg := v.curPkg

						// recursively process this function.
						v.curPkg = pkgIdent.Name
						ast.Walk(v, funcDecl)
						v.curPkg = curPkg
					}
				}
			} else if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				editedFuncName := getEditedFuncName(v.curPkg, ident.Name)
				ident.Name = editedFuncName
				if !v.doneFunc[editedFuncName] {
					funcDecl := v.fv.funcDeclMap[editedFuncName]

					// Add this function in result.
					v.result.Decls = append(v.result.Decls, funcDecl)
					v.doneFunc[editedFuncName] = true

					// recursively process this function.
					ast.Walk(v, funcDecl)
				}
			}
		}
	}

	return v
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Arguments required")
	}

	// TODO: This is hardcoded. Need to fix.
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
	v := visitor{result: res, fv: &fv, doneFunc: make(map[string]bool)}
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

	err = os.MkdirAll("./build", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./build/main.go", goimportedBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
