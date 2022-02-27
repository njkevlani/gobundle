package go_bundle

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/njkevlani/go_bundle/internal/go_bundle/builtinfuncdetector"
	"github.com/njkevlani/go_bundle/internal/go_bundle/collector"
	"github.com/njkevlani/go_bundle/internal/go_bundle/stdpkgdetector"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type visitor struct {
	result   *ast.File
	curPkg   string
	doneFunc map[string]bool
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok && !stdpkgdetector.IsStdPkg(collector.GetFullPkgName(pkgIdent.Name)) && !builtinfuncdetector.IsBuiltinFunc(selectorExpr.Sel.Name) {
					fi := collector.FuncIdentifier{PkgName: pkgIdent.Name, FuncName: selectorExpr.Sel.Name}
					callExpr.Fun = ast.NewIdent(fi.String())

					if !v.doneFunc[fi.String()] {
						funcDecl := collector.GetFuncDecl(fi)

						// Add this function in result.
						v.result.Decls = append(v.result.Decls, funcDecl)
						v.doneFunc[fi.String()] = true
						curPkg := v.curPkg

						// recursively process this function.
						v.curPkg = pkgIdent.Name
						ast.Walk(v, funcDecl)
						v.curPkg = curPkg
					}
				}
			} else if ident, ok := callExpr.Fun.(*ast.Ident); ok && !builtinfuncdetector.IsBuiltinFunc(ident.Name) {
				fi := collector.FuncIdentifier{PkgName: v.curPkg, FuncName: ident.Name}
				ident.Name = fi.String()
				if !v.doneFunc[fi.String()] {
					funcDecl := collector.GetFuncDecl(fi)

					// Add this function in result.
					v.result.Decls = append(v.result.Decls, funcDecl)
					v.doneFunc[fi.String()] = true

					// recursively process this function.
					ast.Walk(v, funcDecl)
				}
			}
		}
	}

	return v
}

func GoBundle(fileName string) ([]byte, error) {
	inFile, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	collector.PutInFileFuncDecls(inFile)

	importPkgs := collector.GetNonStdNonProcessedImports(inFile)

	for len(importPkgs) != 0 {
		pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedSyntax}, importPkgs...)

		if err != nil {
			return nil, err
		}

		collector.PutFuncDecls(pkgs...)

		importPkgs = nil
		for _, pkg := range pkgs {
			for _, file := range pkg.Syntax {
				importPkgs = append(importPkgs, collector.GetNonStdNonProcessedImports(file)...)
			}
		}
	}

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
	v := visitor{result: res, doneFunc: make(map[string]bool)}
	ast.Walk(v, res)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), res)
	if err != nil {
		return nil, err
	}

	return imports.Process("", buf.Bytes(), nil)
}
