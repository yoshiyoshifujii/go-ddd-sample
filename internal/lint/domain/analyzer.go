package domain

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const analyzerName = "domainctor"

var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  "require NewXxx constructor taking at least one param and returning value for exported domain types",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !isDomainPackage(pass) {
		return nil, nil
	}

	exportedTypes := map[string]token.Pos{}
	constructors := map[string]bool{}

	for _, file := range pass.Files {
		if !fileInDomain(pass, file) {
			continue
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !ast.IsExported(typeSpec.Name.Name) {
					continue
				}
				if isInterfaceType(typeSpec.Type) {
					continue
				}
				exportedTypes[typeSpec.Name.Name] = typeSpec.Pos()
			}
		}
	}

	if len(exportedTypes) == 0 {
		return nil, nil
	}

	for _, file := range pass.Files {
		if !fileInDomain(pass, file) {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			name := fn.Name.Name
			if !strings.HasPrefix(name, "New") || name == "New" {
				continue
			}
			typeName := strings.TrimPrefix(name, "New")
			if !hasValueReturn(fn, typeName) || !hasParams(fn) {
				continue
			}
			constructors[typeName] = true
		}
	}

	for name, pos := range exportedTypes {
		if !constructors[name] {
			pass.Reportf(pos, "exported type %s must have constructor New%s taking at least one param and returning %s", name, name, name)
		}
	}

	return nil, nil
}

func isDomainPackage(pass *analysis.Pass) bool {
	for _, file := range pass.Files {
		if fileInDomain(pass, file) {
			return true
		}
	}
	return false
}

func fileInDomain(pass *analysis.Pass, file *ast.File) bool {
	if file == nil {
		return false
	}
	f := pass.Fset.File(file.Pos())
	if f == nil {
		return false
	}
	path := filepath.ToSlash(f.Name())
	return strings.Contains(path, "/internal/domain/")
}

func hasValueReturn(fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
		return false
	}
	if len(fn.Type.Results.List) == 1 {
		return isIdentNamed(fn.Type.Results.List[0].Type, typeName)
	}
	return false
}

func hasParams(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil {
		return false
	}
	return len(fn.Type.Params.List) > 0
}

func isIdentNamed(t ast.Expr, name string) bool {
	ident, ok := t.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == name
}

func isInterfaceType(t ast.Expr) bool {
	_, ok := t.(*ast.InterfaceType)
	return ok
}
