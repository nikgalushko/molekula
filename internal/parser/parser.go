package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"reflect"
	"strings"
)

type File struct {
	Name    string
	Structs map[string][]Field
}

type Field struct {
	Name  string
	Alias string
	Type  ast.Expr
}

func Parse(dir string) (map[string][]File, error) {
	fset := token.NewFileSet()

	pkgs, err := goparser.ParseDir(fset, dir, nil, goparser.ParseComments)
	if err != nil {
		return nil, err
	}

	ret := make(map[string][]File)

	for _, pkg := range pkgs {
		var files []File

		for _, f := range pkg.Files {
			v := &visitor{
				structs: make(map[string][]Field),
			}
			ast.Walk(v, f)

			files = append(files, File{
				Name:    f.Name.Name + ".go",
				Structs: v.structs,
			})
		}

		ret[dir] = files
	}

	return ret, nil
}

type visitor struct {
	structs map[string][]Field
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	s, ok := n.(*ast.TypeSpec)
	if !ok || s == nil {
		return v
	}

	switch t := s.Type.(type) {
	case *ast.StructType:
		description := make([]Field, len(t.Fields.List))

		for i, f := range t.Fields.List {
			description[i] = Field{
				Type:  f.Type,
				Name:  f.Names[0].Name, // TODO:
				Alias: strings.ToLower(f.Names[0].Name),
			}

			if f.Tag != nil {
				tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
				description[i].Alias = tag.Get("molekula")
			}
		}

		v.structs[s.Name.Name] = description
	default:
		fmt.Println(t)
	}

	return v
}
