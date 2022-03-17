package resultmaker

import (
	"go/ast"

	"github.com/njkevlani/gobundle/internal/builtinfuncdetector"
	"github.com/njkevlani/gobundle/internal/collector"
	"github.com/njkevlani/gobundle/internal/stdpkgdetector"
)

type visitor struct {
	result         *ast.File
	curFullPkgName string
	curFilepath    string
	doneDecl       map[string]bool
	localVars      map[string]collector.DeclIdentifier
	ic             *collector.ImportCollector
	dc             *collector.DeclCollector
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			v.handleAssignStmt(assignStmt)
		} else if callExpr, ok := n.(*ast.CallExpr); ok {
			v.handleCallerExpr(callExpr)
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
	} else {
		return
	}

	var di collector.DeclIdentifier
	if unaryExpr, ok := rhs.(*ast.UnaryExpr); ok {
		if compositeLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
			if selectorExpr, ok := compositeLit.Type.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok {
					pkgName := pkgIdent.Name
					structName := selectorExpr.Sel.Name
					di.FullPkgName, di.StructName = v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, pkgName), structName
					compositeLit.Type = ast.NewIdent(v.dc.EditedStructName(di))
				}
			}
		}
	} else if compositeLit, ok := rhs.(*ast.CompositeLit); ok {
		if ident, ok := compositeLit.Type.(*ast.Ident); ok {
			di.FullPkgName, di.StructName = v.curFullPkgName, ident.Name
			compositeLit.Type = ast.NewIdent(v.dc.EditedStructName(di))
		}
	} else if callExpr, ok := rhs.(*ast.CallExpr); ok {
		// Handle call to a function that returns a struct, like NewTrie()
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				fullPkgNameForCalledFunc := v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, ident.Name)
				calledFunc := v.dc.GetDecl(collector.DeclIdentifier{
					FuncName:    selectorExpr.Sel.Name,
					FullPkgName: fullPkgNameForCalledFunc,
				})
				if calledFunc != nil {
					calledFuncCasted := calledFunc.(*ast.FuncDecl)
					if calledFuncCasted.Type.Results != nil && len(calledFuncCasted.Type.Results.List) == 1 {
						if returnType, ok := calledFuncCasted.Type.Results.List[0].Type.(*ast.Ident); ok {
							di.FullPkgName = fullPkgNameForCalledFunc
							di.StructName = returnType.Name
						}
					}
				}
			}
		}
	}

	if funcDecl := v.dc.GetDecl(di); funcDecl != nil {
		v.result.Decls = append(v.result.Decls, funcDecl)
	}
	v.localVars[variableName] = di
}

func (v *visitor) handleCallerExpr(callExpr *ast.CallExpr) {
	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if selectorIdent, ok := selectorExpr.X.(*ast.Ident); ok {
			var di collector.DeclIdentifier

			if builtinfuncdetector.IsBuiltinFunc(selectorExpr.Sel.Name) {
				return
			}

			if diFromV, ok := v.localVars[selectorIdent.Name]; ok {
				di = diFromV
				di.FuncName = selectorExpr.Sel.Name
			} else if fullPkgName := v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, selectorIdent.Name); !stdpkgdetector.IsStdPkg(fullPkgName) {
				di = collector.DeclIdentifier{
					FuncName:    selectorExpr.Sel.Name,
					FullPkgName: fullPkgName,
				}
				callExpr.Fun = ast.NewIdent(v.dc.EditedFuncName(di))
			} else {
				return
			}

			if funcDecl := v.dc.GetDecl(di); funcDecl != nil && !v.doneDecl[di.DeclKey()] {
				v.result.Decls = append(v.result.Decls, funcDecl)

				v.doneDecl[di.DeclKey()] = true

				curFullPkgName, curFilepath, localVars := v.curFullPkgName, v.curFilepath, v.localVars

				// recursively process this function.
				v.curFullPkgName = di.FullPkgName
				v.curFilepath = v.dc.GetDeclFilepath(di)

				v.localVars = make(map[string]collector.DeclIdentifier)
				funcDeclCasted := funcDecl.(*ast.FuncDecl)
				for _, field := range funcDeclCasted.Type.Params.List {
					var funcFieldDi collector.DeclIdentifier
					if starExpr, ok := field.Type.(*ast.StarExpr); ok {
						if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
							funcFieldDi.StructName = selectorExpr.Sel.Name
							if ident, ok := selectorExpr.X.(*ast.Ident); ok {
								funcFieldDi.FullPkgName = v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, ident.Name)
								if !stdpkgdetector.IsStdPkg(funcFieldDi.FullPkgName) {
									starExpr.X = ast.NewIdent(v.dc.EditedStructName(funcFieldDi))
								}
							}
						}
					} else if selectorExpr, ok := field.Type.(*ast.SelectorExpr); ok {
						funcFieldDi.StructName = selectorExpr.Sel.Name
						if ident, ok := selectorExpr.X.(*ast.Ident); ok {
							funcFieldDi.FullPkgName = v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, ident.Name)
						}
						field.Type = ast.NewIdent(v.dc.EditedStructName(funcFieldDi))
					}
					for _, varNameIdent := range field.Names {
						v.localVars[varNameIdent.Name] = funcFieldDi
					}
				}

				ast.Walk(v, funcDecl)

				v.curFullPkgName, v.curFilepath, v.localVars = curFullPkgName, curFilepath, localVars
			}
		}
	} else if ident, ok := callExpr.Fun.(*ast.Ident); ok && !builtinfuncdetector.IsBuiltinFunc(ident.Name) {
		di := collector.DeclIdentifier{
			FuncName:    ident.Name,
			FullPkgName: v.curFullPkgName,
		}

		var funcDecl ast.Decl

		if funcDecl = v.dc.GetDecl(di); funcDecl == nil {
			return
		}

		ident.Name = v.dc.EditedFuncName(di)

		if !v.doneDecl[di.DeclKey()] {
			// Add this function in result.
			v.result.Decls = append(v.result.Decls, funcDecl)
			v.doneDecl[di.DeclKey()] = true

			localVars := v.localVars

			v.localVars = make(map[string]collector.DeclIdentifier)
			funcDeclCasted := funcDecl.(*ast.FuncDecl)
			for _, field := range funcDeclCasted.Type.Params.List {
				var funcFieldDi collector.DeclIdentifier
				if starExpr, ok := field.Type.(*ast.StarExpr); ok {
					if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
						funcFieldDi.StructName = selectorExpr.Sel.Name
						if ident, ok := selectorExpr.X.(*ast.Ident); ok {
							funcFieldDi.FullPkgName = v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, ident.Name)
							if !stdpkgdetector.IsStdPkg(funcFieldDi.FullPkgName) {
								starExpr.X = ast.NewIdent(v.dc.EditedStructName(funcFieldDi))
							}
						}
					}
				} else if selectorExpr, ok := field.Type.(*ast.SelectorExpr); ok {
					funcFieldDi.StructName = selectorExpr.Sel.Name
					if ident, ok := selectorExpr.X.(*ast.Ident); ok {
						funcFieldDi.FullPkgName = v.ic.GetFullPkgName(v.curFullPkgName, v.curFilepath, ident.Name)
					}
					field.Type = ast.NewIdent(v.dc.EditedStructName(funcFieldDi))
				}
				for _, varNameIdent := range field.Names {
					v.localVars[varNameIdent.Name] = funcFieldDi
				}
			}

			// recursively process this function.
			ast.Walk(v, funcDecl)

			v.localVars = localVars
		}
	}
}

func MakeResult(res *ast.File, mainFunc *ast.FuncDecl, ic *collector.ImportCollector, dc *collector.DeclCollector, inFilePkg, inFilepath string) {
	v := visitor{
		result:         res,
		doneDecl:       make(map[string]bool),
		localVars:      make(map[string]collector.DeclIdentifier),
		curFullPkgName: inFilePkg,
		curFilepath:    inFilepath,
		ic:             ic,
		dc:             dc,
	}
	ast.Walk(v, mainFunc)
}
