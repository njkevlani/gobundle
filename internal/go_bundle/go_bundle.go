package go_bundle

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"

	"github.com/njkevlani/go_bundle/internal/go_bundle/builtinfuncdetector"
	"github.com/njkevlani/go_bundle/internal/go_bundle/stdpkgdetector"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var fullPkgNames = make(map[string]string)

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
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && !stdpkgdetector.IsStdPkg(fullPkgNames[pkgIdent.Name]) && !builtinfuncdetector.IsBuiltinFunc(selectorExpr.Sel.Name) {
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
			} else if ident, ok := callExpr.Fun.(*ast.Ident); ok && !builtinfuncdetector.IsBuiltinFunc(ident.Name) {
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

func getImportSpecs(f *ast.File) []*ast.ImportSpec {
	var imports []*ast.ImportSpec

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			break
		}

		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			imports = append(imports, importSpec)
		}
	}

	return imports
}

func GoBundle(fileName string) ([]byte, error) {
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fv := funcDeclCollectorVisitor{funcDeclMap: make(map[string]*ast.FuncDecl)}

	importSpecs := getImportSpecs(inFile)

	for _, importSpec := range importSpecs {
		pkgName := importSpec.Path.Value[1 : len(importSpec.Path.Value)-1]

		if importSpec.Name != nil {
			fullPkgNames[importSpec.Name.Name] = pkgName
		} else {
			pkgNameSplits := strings.Split(pkgName, "/")
			fullPkgNames[pkgNameSplits[len(pkgNameSplits)-1]] = pkgName
		}
		if !stdpkgdetector.IsStdPkg(pkgName) {
			pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedSyntax}, pkgName)

			if err != nil {
				return nil, err
			}
			for _, pkg := range pkgs {
				for _, file := range pkg.Syntax {
					fv.curPkg = file.Name.Name
					ast.Walk(fv, file)
				}
			}
		}
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
		return nil, errors.New("No main function found in file: " + fileName)
	}

	res := &ast.File{Name: ast.NewIdent("main")}
	res.Decls = append(res.Decls, mainFunc)
	v := visitor{result: res, fv: &fv, doneFunc: make(map[string]bool)}
	ast.Walk(v, res)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), res)
	if err != nil {
		return nil, err
	}

	bytes := buf.Bytes()
	goimportedBytes, err := imports.Process("", bytes, nil)
	if err != nil {
		return nil, err
	}

	goimportedBytes = append([]byte("// Auto generated using https://github.com/njkevlani/go_bundle\n"), goimportedBytes...)

	return goimportedBytes, nil
}
