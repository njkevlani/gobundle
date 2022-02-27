package resultmaker

import (
	"go/ast"

	"github.com/njkevlani/go_bundle/internal/builtinfuncdetector"
	"github.com/njkevlani/go_bundle/internal/collector"
	"github.com/njkevlani/go_bundle/internal/stdpkgdetector"
)

type visitor struct {
	result   *ast.File
	curPkg   string
	doneFunc map[string]bool
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			v.handleCallerExpr(callExpr)
		}
	}

	return v
}

func (v *visitor) handleCallerExpr(callExpr *ast.CallExpr) {
	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok &&
			!stdpkgdetector.IsStdPkg(collector.GetFullPkgName(pkgIdent.Name)) &&
			!builtinfuncdetector.IsBuiltinFunc(selectorExpr.Sel.Name) {
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

func MakeResult(res *ast.File) {
	v := visitor{result: res, doneFunc: make(map[string]bool)}
	ast.Walk(v, res)
}
