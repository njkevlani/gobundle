package collector

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/njkevlani/gobundle/internal/stdpkgdetector"
)

type ImportCollector struct {
	fullPkgNames    map[string]map[string]map[string]string
	importProcessed map[string]bool
}

func NewImportCollector() *ImportCollector {
	return &ImportCollector{
		fullPkgNames:    make(map[string]map[string]map[string]string),
		importProcessed: make(map[string]bool),
	}
}

func (ic *ImportCollector) GetNonStdNonProcessedImports(f *ast.File, fullPkgName, filepath string) []string {
	var imports []string

	if _, exists := ic.fullPkgNames[fullPkgName]; !exists {
		ic.fullPkgNames[fullPkgName] = make(map[string]map[string]string)
	}

	if _, exists := ic.fullPkgNames[fullPkgName][filepath]; !exists {
		ic.fullPkgNames[fullPkgName][filepath] = make(map[string]string)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			break
		}

		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)

			fullImportPath := importSpec.Path.Value[1 : len(importSpec.Path.Value)-1]

			if importSpec.Name != nil {
				ic.fullPkgNames[fullPkgName][filepath][importSpec.Name.Name] = fullImportPath
			} else {
				pkgNameSplits := strings.Split(fullImportPath, "/")

				shortPkgName := pkgNameSplits[len(pkgNameSplits)-1]

				// TODO: This is not alway true.
				// There are cases where this fails.
				// eg: github.com/mattn/go-isatty and default package name is isatty
				// Doing this hack for now, need to understand
				//   1. How does this work?
				//   2. What is proper way to do this?
				//   3. If there are any other such similar cases.
				if strings.HasPrefix(shortPkgName, "go-") {
					shortPkgName = shortPkgName[3:]
				}

				ic.fullPkgNames[fullPkgName][filepath][shortPkgName] = fullImportPath
			}

			if !stdpkgdetector.IsStdPkg(fullImportPath) && !ic.importProcessed[fullImportPath] {
				imports = append(imports, fullImportPath)
				ic.importProcessed[fullImportPath] = true
			}
		}
	}

	return imports
}

func (ic *ImportCollector) GetFullPkgName(pkgName, filepath, shortPkgName string) string {
	if ic.fullPkgNames[pkgName] == nil ||
		ic.fullPkgNames[pkgName][filepath] == nil ||
		ic.fullPkgNames[pkgName][filepath][shortPkgName] == "" {
		panic(fmt.Sprintf(
			"Could not get fullPkgName. pkgName=%s filepath=%s shortPkgName=%s\n", pkgName, filepath, shortPkgName,
		))
	}

	return ic.fullPkgNames[pkgName][filepath][shortPkgName]
}

func (ic *ImportCollector) Debug() {
	fmt.Println("fullPkgNames:")
	spew.Dump(ic.fullPkgNames)
}
