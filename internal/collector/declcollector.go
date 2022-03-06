package collector

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/davecgh/go-spew/spew"
)

type declType int

const (
	funcType declType = iota
	funcOnStructType
	structType
)

type DeclIdentifier struct {
	FuncName    string
	StructName  string
	FullPkgName string
	Filepath    string
}

type DeclCollector struct {
	decls       map[string]ast.Decl
	editedNames map[string]string
	usedNames   map[string]bool
	declFileMap map[string]string
}

func NewDeclCollector() *DeclCollector {
	return &DeclCollector{
		decls:       make(map[string]ast.Decl),
		editedNames: make(map[string]string),
		usedNames:   make(map[string]bool),
		declFileMap: make(map[string]string),
	}
}

func (dc *DeclCollector) EditedFuncName(fi DeclIdentifier) string {
	return dc.editedNames[fi.DeclKey()]
}

func (di DeclIdentifier) DeclKey() string {
	return fmt.Sprintf("%s_%s_%s", di.FullPkgName, di.StructName, di.FuncName)
}

func (dc *DeclCollector) putEditedFuncName(fi DeclIdentifier, dt declType) {
	switch dt {
	case funcOnStructType:
		dc.editedNames[fi.DeclKey()] = fi.FuncName
	case funcType:
		funcName := fi.FuncName
		i := 1
		for dc.usedNames[funcName] {
			i++
			funcName = fmt.Sprintf("%s%d", fi.FuncName, i)
		}

		dc.usedNames[funcName] = true
		dc.editedNames[fi.DeclKey()] = funcName
	case structType:
		structName := fi.StructName
		i := 1
		for dc.usedNames[structName] {
			i++
			structName = fmt.Sprintf("%s%d", fi.StructName, i)
		}

		dc.usedNames[structName] = true
		dc.editedNames[fi.DeclKey()] = structName
	}
}

func (dc *DeclCollector) EditedStructName(fi DeclIdentifier) string {
	return dc.editedNames[fi.DeclKey()]
}

type funcDeclCollectorVisitor struct {
	curFullPkgName string
	curFilepath    string
	declMap        map[string]ast.Decl
	c              *DeclCollector
}

func (fv funcDeclCollectorVisitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			dt := funcType
			di := DeclIdentifier{FuncName: funcDecl.Name.Name, FullPkgName: fv.curFullPkgName, Filepath: fv.curFilepath}
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) == 1 {
				receiver := funcDecl.Recv.List[0]
				if starExpr, ok := receiver.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						di.StructName = ident.Name
						dt = funcOnStructType
					}
				} else if ident, ok := receiver.Type.(*ast.Ident); ok {
					di.StructName = ident.Name
					dt = funcOnStructType
				}
			}

			fv.c.decls[di.DeclKey()] = funcDecl
			fv.c.declFileMap[di.DeclKey()] = fv.curFilepath
			fv.c.putEditedFuncName(di, dt)
		} else if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			di := DeclIdentifier{
				StructName:  genDecl.Specs[0].(*ast.TypeSpec).Name.Name,
				FullPkgName: fv.curFullPkgName,
				Filepath:    fv.curFilepath,
			}

			fv.c.decls[di.DeclKey()] = genDecl
			fv.c.declFileMap[di.DeclKey()] = fv.curFilepath
			fv.c.putEditedFuncName(di, structType)
		}
	}

	return fv
}

func (dc *DeclCollector) CollectFileDecls(file *ast.File, fullPkgName, filepath string) {
	fv := funcDeclCollectorVisitor{
		curFullPkgName: fullPkgName,
		curFilepath:    filepath,
		c:              dc,
	}
	ast.Walk(fv, file)
}

func (dc *DeclCollector) GetDecl(di DeclIdentifier) ast.Decl {
	return dc.decls[di.DeclKey()]
}

func (dc *DeclCollector) GetDeclFilepath(di DeclIdentifier) string {
	return dc.declFileMap[di.DeclKey()]
}

func (dc *DeclCollector) Debug() {
	fmt.Println("editedNames:")
	spew.Dump(dc.editedNames)

	fmt.Println("declFileMap:")
	spew.Dump(dc.declFileMap)

	fmt.Println("decl keys:")
	for key := range dc.decls {
		fmt.Println(key)
	}
}
