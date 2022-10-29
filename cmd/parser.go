package main

import (
	_ "github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"

	"fmt"
	"go/ast"
	"go/build"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/loader"
)

// A PackageInformation contains all the information related to a parsed package.
type PackageInformation struct {
	Name  string
	FQN   string
	files []*ast.File
	defs  map[*ast.Ident]types.Object
}

// parsePackage parses the package in the given directory and returns it.
func parsePackage(packageName string) (*PackageInformation, error) {
	p, err := build.Import(packageName, "", build.FindOnly)
	if err != nil {
		return nil, fmt.Errorf("provided directory (%s) may not under GOPATH (%s): %v",
			packageName, build.Default.GOPATH, err)
	}

	conf := loader.Config{TypeChecker: types.Config{FakeImportC: true}}
	conf.Import(p.ImportPath)
	program, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("couldn't load package: %v", err)
	}

	pkgInfo := program.Package(p.ImportPath)
	return &PackageInformation{
		Name:  pkgInfo.Pkg.Name(),
		FQN:   pkgInfo.Pkg.Path(),
		files: pkgInfo.Files,
		defs:  pkgInfo.Defs,
	}, nil
}

// ValuesOfType generate produces the String method for the named type.
func (pkg *PackageInformation) ValuesOfType(typeName string) ([]string, error) {
	var values, inspectErrs []string
	for _, file := range pkg.files {
		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				// We only care about const declarations.
				return true
			}

			typeSpec := decl.Specs[0].(*ast.TypeSpec)
			if typeName == typeSpec.Name.String() {
				values = append(values, typeSpec.Name.String())
			}
			//if vs, err := pkg.valuesOfTypeIn(typeName, decl); err != nil {
			//	inspectErrs = append(inspectErrs, err.Error())
			//} else {
			//	values = append(values, vs...)
			//}
			return false
		})
	}
	if len(inspectErrs) > 0 {
		return nil, fmt.Errorf("inspecting code:\n\t%v", strings.Join(inspectErrs, "\n\t"))
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no values defined for type %s", typeName)
	}
	return values, nil
}

func (pkg *PackageInformation) findFunctionsAndParameters(typeName string) (map[string][]*ast.Ident, error) {
	result := map[string][]*ast.Ident{}
	for _, file := range pkg.files {
		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok {
				// We only care about const declarations.
				return true
			}

			if funcDecl.Recv != nil {
				if funcDecl.Name.Name[:1] == strings.ToUpper(funcDecl.Name.Name[:1]) {
					for _, recv := range funcDecl.Recv.List {
						typ := recv.Type
						starExp, ok := typ.(*ast.StarExpr)
						if ok {
							typ = starExp.X
						}
						if ident, ok := typ.(*ast.Ident); ok {
							if ident.Name == typeName {
								var pTypes []*ast.Ident
								if funcDecl.Type.Params != nil {
									for _, param := range funcDecl.Type.Params.List {
										p, ok := param.Type.(*ast.Ident)
										if ok {
											pTypes = append(pTypes, p)
										} else {
											log.Printf(fmt.Sprintf("unknown identifier, func %s(?) : %v", funcDecl.Name, param.Type))
										}
									}
								}
								result[funcDecl.Name.String()] = pTypes
							}
						}
					}
				}
			}
			return false
		})
	}
	return result, nil
}

func (pkg *PackageInformation) valuesOfTypeIn(typeName string, decl *ast.GenDecl) ([]string, error) {
	var values []string

	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value, the constant is untyped.
			// Skip this vspec and reset the remembered type.
			typ = ""
			continue
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != typeName {
			// This is not the type we're looking for.
			continue
		}

		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := pkg.defs[name]
			if !ok {
				return nil, fmt.Errorf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				return nil, fmt.Errorf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			values = append(values, name.Name)
		}
	}
	return values, nil
}
