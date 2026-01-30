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
	Doc:  "require NewXxx constructor returning value for exported domain structs",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !isDomainPackage(pass) {
		return nil, nil
	}

	exportedStructs := map[string]token.Pos{}
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
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					continue
				}
				exportedStructs[typeSpec.Name.Name] = typeSpec.Pos()
			}
		}
	}

	if len(exportedStructs) == 0 {
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
			if !hasValueReturn(fn, typeName) {
				continue
			}
			constructors[typeName] = true
		}
	}

	for name, pos := range exportedStructs {
		if !constructors[name] {
			pass.Reportf(pos, "exported struct %s must have constructor New%s returning %s", name, name, name)
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
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}
	result := fn.Type.Results.List[0].Type
	ident, ok := result.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == typeName
}
