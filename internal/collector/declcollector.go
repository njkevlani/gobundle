package collector

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

var funcDecls = make(map[string]*ast.FuncDecl)
var genDecls = make(map[string]*ast.GenDecl)

type FuncIdentifier struct {
	FuncName   string
	StructName string
	PkgName    string
}

func (fi FuncIdentifier) EditedFuncName() string {
	if fi.StructName != "" {
		return fi.FuncName
	}

	if fi.PkgName == "" {
		return fi.FuncName
	}

	return fmt.Sprintf("%s_%s", fi.PkgName, fi.FuncName)
}

func (fi FuncIdentifier) EditedStructName() string {
	if fi.PkgName == "" {
		return fi.StructName
	}

	return fmt.Sprintf("%s_%s", fi.PkgName, fi.StructName)
}

func (fi FuncIdentifier) DeclKey() string {
	return fi.PkgName + "_" + fi.StructName + "_" + fi.FuncName
}

type funcDeclCollectorVisitor struct {
	curPkg      string
	funcDeclMap map[string]*ast.FuncDecl
	genDeclMap  map[string]*ast.GenDecl
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			fi := FuncIdentifier{PkgName: fv.curPkg, FuncName: funcDecl.Name.Name}
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) == 1 {
				receiver := funcDecl.Recv.List[0]
				if starExpr, ok := receiver.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						fi.StructName = ident.Name
						ident.Name = fi.EditedStructName()
					}
				} else if ident, ok := receiver.Type.(*ast.Ident); ok {
					fi.StructName = ident.Name
					ident.Name = fi.EditedStructName()
				}
			}

			funcDecl.Name.Name = fi.EditedFuncName()

			fv.funcDeclMap[fi.DeclKey()] = funcDecl
		} else if genDecl, ok := n.(*ast.GenDecl); ok && len(genDecl.Specs) == 1 {
			if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
				fi := FuncIdentifier{PkgName: fv.curPkg, StructName: typeSpec.Name.Name}
				typeSpec.Name.Name = fi.EditedStructName()

				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if starExpr, ok := field.Type.(*ast.StarExpr); ok {
							// TODO: This will only work for renaming self reference in struct member.
							if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == fi.StructName {
								ident.Name = fi.EditedStructName()
							}
						}
					}
				}

				// No need to recursively get all the structs because we will get all the things from this import + recursive imports
				// But this is not great, because it will not work when struct a depends on struct b and we load struct a before b

				fv.genDeclMap[fi.DeclKey()] = genDecl
			}
		} else if assignStmt, ok := n.(*ast.AssignStmt); ok {
			for _, expr := range assignStmt.Rhs {
				if unaryExpr, ok := expr.(*ast.UnaryExpr); ok {
					if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
						if ident, ok := compositeLit.Type.(*ast.Ident); ok {
							fi := FuncIdentifier{PkgName: fv.curPkg, StructName: ident.Name}

							if _, ok := fv.genDeclMap[fi.DeclKey()]; ok {
								ident.Name = fi.EditedStructName()
							}
						}
					}
				}

			}
		}
	}

	return fv
}

func PutFuncDecls(pkgs ...*packages.Package) {
	fv := funcDeclCollectorVisitor{funcDeclMap: funcDecls, genDeclMap: genDecls}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fv.curPkg = file.Name.Name
			ast.Walk(fv, file)
		}
	}
}

func PutInFileFuncDecls(inFile *ast.File) {
	fv := funcDeclCollectorVisitor{funcDeclMap: funcDecls, genDeclMap: genDecls}
	ast.Walk(fv, inFile)
}

func GetFuncDecl(fi FuncIdentifier) *ast.FuncDecl {
	return funcDecls[fi.DeclKey()]
}

func GetGenDecl(fi FuncIdentifier) *ast.GenDecl {
	return genDecls[fi.DeclKey()]
}
