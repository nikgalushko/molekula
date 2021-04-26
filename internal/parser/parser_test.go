package parser

import (
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/nikgalushko/molekula/internal/ast"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	fset := token.NewFileSet()

	pkgs, err := goparser.ParseDir(fset, "testdata/", nil, goparser.ParseComments)
	assert.NoError(t, err)

	var objects []Object
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			objects = Parse(f)
		}
	}

	assert.Equal(t, Object{
		BinName: "data",
		Type: ast.Map{
			Key: ast.BuiltIn("string"),
			Value: ast.Map{
				Key:   ast.BuiltIn("string"),
				Value: ast.BuiltIn("int"),
			},
		},
	}, find(objects, "data"))

	assert.Equal(t, Object{
		BinName: "kek",
		Type: ast.Struct{
			Name: "Foo",
			Fields: []ast.StructField{
				{
					Name:  "Str",
					Alias: "version",
					Type:  ast.BuiltIn("string"),
				},
				{
					Name:  "Intrf",
					Alias: "intrf",
					Type:  ast.BuiltIn("interface{}"),
				},
				{
					Name:  "ArrInt",
					Alias: "arrint",
					Type:  ast.Array{Element: ast.BuiltIn("int64")},
				},
				{
					Name:  "Float",
					Alias: "float",
					Type:  ast.BuiltIn("float64"),
				},
				{
					Name:  "Int",
					Alias: "int",
					Type:  ast.BuiltIn("int64"),
				},
			},
		},
	}, find(objects, "kek"))

	assert.Equal(t, Object{
		BinName: "config_version",
		Type:    ast.BuiltIn("int"),
	}, find(objects, "config_version"))

	assert.Equal(t, Object{
		BinName: "weights",
		Type:    ast.Array{Element: ast.BuiltIn("float64")},
	}, find(objects, "weights"))

	assert.Equal(t, Object{
		BinName: "config",
		Type: ast.Map{
			Key: ast.BuiltIn("string"),
			Value: ast.Struct{
				Name: "Value",
				Fields: []ast.StructField{
					{Name: "Gender", Alias: "gender", Type: ast.BuiltIn("string")},
					{Name: "ID", Alias: "id", Type: ast.BuiltIn("int64")},
				},
			},
		},
	}, find(objects, "config"))

	assert.Equal(t, Object{
		BinName: "slice",
		Type: ast.Array{
			Element: ast.Struct{
				Name: "Value",
				Fields: []ast.StructField{
					{Name: "Gender", Alias: "gender", Type: ast.BuiltIn("string")},
					{Name: "ID", Alias: "id", Type: ast.BuiltIn("int64")},
				},
			},
		},
	}, find(objects, "slice"))
}

func find(objects []Object, name string) Object {
	for _, o := range objects {
		if o.BinName == name {
			return o
		}
	}

	return Object{}
}
