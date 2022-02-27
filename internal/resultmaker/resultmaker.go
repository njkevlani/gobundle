package resultmaker

import (
	"go/ast"

	"github.com/njkevlani/go_bundle/internal/builtinfuncdetector"
	"github.com/njkevlani/go_bundle/internal/collector"
	"github.com/njkevlani/go_bundle/internal/stdpkgdetector"
)

type visitor struct {
	result          *ast.File
	curPkg          string
	doneFunc        map[string]bool
	localStructVars map[string]collector.FuncIdentifier
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			v.handleCallerExpr(callExpr)
		} else if assignStmt, ok := n.(*ast.AssignStmt); ok {
			v.handleAssignStmt(assignStmt)
		}
	}

	return v
}

func (v *visitor) handleAssignStmt(assignStmt *ast.AssignStmt) {
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return
	}

	lhs, rhs := assignStmt.Lhs[0], assignStmt.Rhs[0]

	var variableName string

	if ident, ok := lhs.(*ast.Ident); ok {
		variableName = ident.Name
	}

	fi := collector.FuncIdentifier{}
	if unaryExpr, ok := rhs.(*ast.UnaryExpr); ok {
		if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
			if selectorExpr, ok := compositeLit.Type.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok {
					pkgName := pkgIdent.Name
					structName := selectorExpr.Sel.Name
					fi = collector.FuncIdentifier{PkgName: pkgName, StructName: structName}
					compositeLit.Type = ast.NewIdent(fi.EditedStructName())
				}
			}
		}
	} else if compositeLit, ok := rhs.(*ast.CompositeLit); ok {
		if ident, ok := compositeLit.Type.(*ast.Ident); ok {
			fi = collector.FuncIdentifier{PkgName: v.curPkg, StructName: ident.Name}
			compositeLit.Type = ast.NewIdent(fi.EditedStructName())
		}
	}

	if fi.StructName != "" {
		v.result.Decls = append(v.result.Decls, collector.GetGenDecl(fi))
		v.localStructVars[variableName] = fi
	}
}

func (v *visitor) handleCallerExpr(callExpr *ast.CallExpr) {
	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok &&
			!stdpkgdetector.IsStdPkg(collector.GetFullPkgName(pkgIdent.Name)) &&
			!builtinfuncdetector.IsBuiltinFunc(selectorExpr.Sel.Name) {

			var fi collector.FuncIdentifier

			if fiFromV, ok := v.localStructVars[pkgIdent.Name]; ok {
				fi = fiFromV
				fi.FuncName = selectorExpr.Sel.Name
			} else {
				fi = collector.FuncIdentifier{PkgName: pkgIdent.Name, FuncName: selectorExpr.Sel.Name}
				callExpr.Fun = ast.NewIdent(fi.EditedFuncName())
			}

			if !v.doneFunc[fi.DeclKey()] {
				funcDecl := collector.GetFuncDecl(fi)

				// Add this function in result.
				v.result.Decls = append(v.result.Decls, funcDecl)
				v.doneFunc[fi.DeclKey()] = true
				curPkg := v.curPkg
				localStructVars := v.localStructVars

				// recursively process this function.
				v.curPkg = pkgIdent.Name
				v.localStructVars = make(map[string]collector.FuncIdentifier)

				ast.Walk(v, funcDecl)

				v.curPkg = curPkg
				v.localStructVars = localStructVars
			}
		}
	} else if ident, ok := callExpr.Fun.(*ast.Ident); ok && !builtinfuncdetector.IsBuiltinFunc(ident.Name) {
		fi := collector.FuncIdentifier{PkgName: v.curPkg, FuncName: ident.Name}
		ident.Name = fi.EditedFuncName()
		if !v.doneFunc[fi.DeclKey()] {
			funcDecl := collector.GetFuncDecl(fi)

			// Add this function in result.
			v.result.Decls = append(v.result.Decls, funcDecl)
			v.doneFunc[fi.DeclKey()] = true

			// recursively process this function.
			ast.Walk(v, funcDecl)
		}
	}
}

func MakeResult(res *ast.File) {
	v := visitor{result: res, doneFunc: make(map[string]bool), localStructVars: make(map[string]collector.FuncIdentifier)}
	ast.Walk(v, res)
}
