package go_bundle

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/njkevlani/go_bundle/internal/collector"
	"github.com/njkevlani/go_bundle/internal/resultmaker"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

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

	resultmaker.MakeResult(res)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), res)
	if err != nil {
		return nil, err
	}

	return imports.Process("", buf.Bytes(), nil)
}
