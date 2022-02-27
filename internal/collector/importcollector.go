package collector

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/njkevlani/go_bundle/internal/stdpkgdetector"
)

var fullPkgNames = make(map[string]string)
var importProcessed = make(map[string]bool)

func GetNonStdNonProcessedImports(f *ast.File) []string {
	var imports []string

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			break
		}

		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)

			fullImportPath := importSpec.Path.Value[1 : len(importSpec.Path.Value)-1]

			// TODO: fullPkgName depends on which file we are talking about.
			if importSpec.Name != nil {
				fullPkgNames[importSpec.Name.Name] = fullImportPath
			} else {
				pkgNameSplits := strings.Split(fullImportPath, "/")
				fullPkgNames[pkgNameSplits[len(pkgNameSplits)-1]] = fullImportPath
			}

			if !stdpkgdetector.IsStdPkg(fullImportPath) && !importProcessed[fullImportPath] {
				imports = append(imports, fullImportPath)
				importProcessed[fullImportPath] = true
			}
		}
	}

	return imports
}

func GetFullPkgName(pkgName string) string {
	return fullPkgNames[pkgName]
}
