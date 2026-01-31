package usecase

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const analyzerName = "usecaseexec"

var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  "require Execute(ctx, XxxUsecaseInput) (*XxxUsecaseOutput, error) for exported usecase types",
	Run:  run,
}

type execState struct {
	pos        token.Pos
	hasExecute bool
	valid      bool
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !isUsecasePackage(pass) {
		return nil, nil
	}

	usecaseTypes := map[string]*execState{}

	for _, file := range pass.Files {
		if !fileInUsecase(pass, file) {
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
				if isUsecaseDTO(typeSpec.Name.Name) {
					continue
				}
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					continue
				}
				usecaseTypes[typeSpec.Name.Name] = &execState{pos: typeSpec.Pos()}
			}
		}
	}

	if len(usecaseTypes) == 0 {
		return nil, nil
	}

	for _, file := range pass.Files {
		if !fileInUsecase(pass, file) {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name == nil {
				continue
			}
			if fn.Name.Name != "Execute" {
				continue
			}
			recvTypeName := receiverTypeName(fn.Recv)
			if recvTypeName == "" {
				continue
			}
			state, ok := usecaseTypes[recvTypeName]
			if !ok {
				continue
			}
			state.hasExecute = true
			if hasExecuteSignature(fn, recvTypeName) {
				state.valid = true
				continue
			}
			pass.Reportf(fn.Pos(), "Execute must have signature Execute(ctx context.Context, %sUsecaseInput) (*%sUsecaseOutput, error)", recvTypeName, recvTypeName)
		}
	}

	for name, state := range usecaseTypes {
		if !state.hasExecute {
			pass.Reportf(state.pos, "exported type %s must have method Execute(ctx context.Context, %sUsecaseInput) (*%sUsecaseOutput, error)", name, name, name)
		}
	}

	return nil, nil
}

func isUsecasePackage(pass *analysis.Pass) bool {
	for _, file := range pass.Files {
		if fileInUsecase(pass, file) {
			return true
		}
	}
	return false
}

func isUsecaseDTO(name string) bool {
	return strings.HasSuffix(name, "Input") ||
		strings.HasSuffix(name, "Output") ||
		strings.HasSuffix(name, "UsecaseInput") ||
		strings.HasSuffix(name, "UsecaseOutput")
}

func fileInUsecase(pass *analysis.Pass, file *ast.File) bool {
	if file == nil {
		return false
	}
	f := pass.Fset.File(file.Pos())
	if f == nil {
		return false
	}
	path := filepath.ToSlash(f.Name())
	return strings.Contains(path, "/internal/usecase/")
}

func receiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	switch t := recv.List[0].Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

func hasExecuteSignature(fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Params == nil || fn.Type.Results == nil {
		return false
	}
	if len(fn.Type.Params.List) != 2 {
		return false
	}
	if len(fn.Type.Results.List) != 2 {
		return false
	}

	if !isSelectorNamed(fn.Type.Params.List[0].Type, "context", "Context") {
		return false
	}
	if !isIdentNamed(fn.Type.Params.List[1].Type, typeName+"UsecaseInput") {
		return false
	}
	if !isPointerToIdentNamed(fn.Type.Results.List[0].Type, typeName+"UsecaseOutput") {
		return false
	}
	if !isIdentNamed(fn.Type.Results.List[1].Type, "error") {
		return false
	}
	return true
}

func isSelectorNamed(t ast.Expr, pkg, name string) bool {
	sel, ok := t.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == pkg && sel.Sel.Name == name
}

func isIdentNamed(t ast.Expr, name string) bool {
	ident, ok := t.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == name
}

func isPointerToIdentNamed(t ast.Expr, name string) bool {
	star, ok := t.(*ast.StarExpr)
	if !ok {
		return false
	}
	return isIdentNamed(star.X, name)
}
