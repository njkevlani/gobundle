package collector

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

var funcDelcs = make(map[string]*ast.FuncDecl)

type FuncIdentifier struct {
	FuncName string
	PkgName  string
}

func (fi FuncIdentifier) String() string {
	if fi.PkgName == "" {
		return fi.FuncName
	}

	return fmt.Sprintf("%s_%s", fi.PkgName, fi.FuncName)
}

type funcDeclCollectorVisitor struct {
	curPkg      string
	funcDeclMap map[string]*ast.FuncDecl
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			funcDecl.Name.Name = FuncIdentifier{PkgName: fv.curPkg, FuncName: funcDecl.Name.Name}.String()
			fv.funcDeclMap[funcDecl.Name.Name] = funcDecl
		}
	}

	return fv
}

func PutFuncDecls(pkgs ...*packages.Package) {
	fv := funcDeclCollectorVisitor{funcDeclMap: funcDelcs}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fv.curPkg = file.Name.Name
			ast.Walk(fv, file)
		}
	}
}

func PutInFileFuncDecls(inFile *ast.File) {
	fv := funcDeclCollectorVisitor{funcDeclMap: funcDelcs}
	ast.Walk(fv, inFile)
}

func GetFuncDecl(fi FuncIdentifier) *ast.FuncDecl {
	return funcDelcs[fi.String()]
}
