package gobundle

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/njkevlani/gobundle/internal/collector"
	"github.com/njkevlani/gobundle/internal/resultmaker"
)

func GoBundle(filepath string) ([]byte, error) {
	fileset := token.NewFileSet()
	inFile, err := parser.ParseFile(fileset, filepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	importCollector := collector.NewImportCollector()

	declCollector := collector.NewDeclCollector()

	declCollector.CollectFileDecls(inFile, inFile.Name.Name, filepath)

	importPkgs := importCollector.GetNonStdNonProcessedImports(inFile, inFile.Name.Name, filepath)

	for len(importPkgs) != 0 {
		pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, importPkgs...)

		if err != nil {
			return nil, err
		}

		importPkgs = nil
		for _, pkg := range pkgs {
			for _, pkgFilepath := range pkg.GoFiles {
				file, err := parser.ParseFile(fileset, pkgFilepath, nil, parser.ParseComments)
				if err != nil {
					return nil, err
				}

				declCollector.CollectFileDecls(file, pkg.ID, pkgFilepath)

				importPkgs = append(importPkgs, importCollector.GetNonStdNonProcessedImports(file, pkg.ID, pkgFilepath)...)
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
		return nil, errors.New("No main function found in file: " + filepath)
	}

	res := &ast.File{Name: ast.NewIdent("main")}
	res.Decls = append(res.Decls, mainFunc)

	resultmaker.MakeResult(res, mainFunc, importCollector, declCollector, inFile.Name.Name, filepath)

	printConfig := &printer.Config{Mode: printer.TabIndent, Tabwidth: 1}

	var buf bytes.Buffer
	err = printConfig.Fprint(&buf, token.NewFileSet(), res)
	if err != nil {
		return nil, err
	}

	return imports.Process("", buf.Bytes(), nil)
	// return []byte{}, nil
}
