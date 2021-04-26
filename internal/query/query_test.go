package query

import (
	"testing"

	"github.com/nikgalushko/molekula/internal/ast"
	"github.com/nikgalushko/molekula/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestBuild_SimpleArray(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Array{
			Element: ast.BuiltIn("int"),
		},
	})

	assert.Equal(t, Query{
		IsTop: true, IsArray: true,
		Type: "[]int",
		Next: &Query{IsBuiltin: true, Index: 1, Type: "int"},
	}, q)
}

func TestBuild_NestedSimpleArray(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Array{
			Element: ast.Array{
				Element: ast.BuiltIn("int"),
			},
		},
	})

	assert.Equal(t, Query{
		IsTop: true, IsArray: true,
		Type: "[][]int",
		Next: &Query{
			IsArray: true,
			Index:   1,
			Type:    "[]int",
			Next:    &Query{IsBuiltin: true, Index: 2, Type: "int"},
		},
	}, q)
}

func TestBuild_ArrayOfMap(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Array{
			Element: ast.Map{
				Key:   ast.BuiltIn("int"),
				Value: ast.BuiltIn("string"),
			},
		},
	})

	assert.Equal(t, Query{
		IsTop: true, IsArray: true,
		Type: "[]map[int]string",
		Next: &Query{
			IsMap:   true,
			Index:   1,
			Type:    "map[int]string",
			KeyType: "int",
			Next:    &Query{IsBuiltin: true, Index: 2, Type: "string"},
		},
	}, q)
}

func TestBuild_SimpleMap(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Map{
			Key:   ast.BuiltIn("int"),
			Value: ast.BuiltIn("string"),
		},
	})

	assert.Equal(t, Query{
		IsTop: true, IsMap: true,
		Index:   0,
		Type:    "map[int]string",
		KeyType: "int",
		Next:    &Query{IsBuiltin: true, Index: 1, Type: "string"},
	}, q)
}

func TestBuild_MapOfMap(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Map{
			Key: ast.BuiltIn("int"),
			Value: ast.Map{
				Key: ast.BuiltIn("string"),
				Value: ast.Array{
					Element: ast.BuiltIn("float64"),
				},
			},
		},
	})

	assert.Equal(t, Query{
		IsTop: true, IsMap: true,
		Type:    "map[int]map[string][]float64",
		KeyType: "int",
		Next: &Query{
			IsMap:   true,
			Index:   1,
			Type:    "map[string][]float64",
			KeyType: "string",
			Next: &Query{
				IsArray: true,
				Index:   2,
				Type:    "[]float64",
				Next:    &Query{IsBuiltin: true, Index: 3, Type: "float64"},
			},
		},
	}, q)
}

func TestBuild_Struct(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Struct{
			Name: "Foo",
			Fields: []ast.StructField{
				{
					Name:  "Gender",
					Alias: "gender",
					Type:  ast.BuiltIn("string"),
				},
				{
					Name:  "ID",
					Alias: "id",
					Type:  ast.BuiltIn("int64"),
				},
			},
		},
	})

	assert.Equal(t, Query{
		IsTop:    true,
		IsStruct: true,
		Type:     "Foo",
		Fields: []Query{
			{Name: "Gender", Alias: "gender", Type: "string", Index: 1, IsBuiltin: true},
			{Name: "ID", Alias: "id", Type: "int64", Index: 1, IsBuiltin: true},
		},
	}, q)
}

func TestBuild_MapOfStruct(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Map{
			Key: ast.BuiltIn("int"),
			Value: ast.Struct{
				Name: "Foo",
				Fields: []ast.StructField{
					{
						Name:  "Gender",
						Alias: "gender",
						Type:  ast.BuiltIn("string"),
					},
					{
						Name:  "ID",
						Alias: "id",
						Type:  ast.BuiltIn("int64"),
					},
				},
			},
		},
	})

	assert.Equal(t, Query{
		IsTop:   true,
		IsMap:   true,
		Type:    "map[int]Foo",
		KeyType: "int",
		Next: &Query{
			IsStruct: true,
			Index:    1,
			Type:     "Foo",
			Fields: []Query{
				{Name: "Gender", Alias: "gender", Type: "string", Index: 2, IsBuiltin: true},
				{Name: "ID", Alias: "id", Type: "int64", Index: 2, IsBuiltin: true},
			},
		},
	}, q)
}

func TestBuild_ArrayOfStruct(t *testing.T) {
	q := Build(parser.Object{
		Type: ast.Array{
			Element: ast.Struct{
				Name: "Foo",
				Fields: []ast.StructField{
					{
						Name:  "Gender",
						Alias: "gender",
						Type:  ast.BuiltIn("string"),
					},
					{
						Name:  "ID",
						Alias: "id",
						Type:  ast.BuiltIn("int64"),
					},
				},
			},
		},
	})

	assert.Equal(t, Query{
		IsTop:   true,
		IsArray: true,
		Type:    "[]Foo",
		Next: &Query{
			IsStruct: true,
			Index:    1,
			Type:     "Foo",
			Fields: []Query{
				{Name: "Gender", Alias: "gender", Type: "string", Index: 2, IsBuiltin: true},
				{Name: "ID", Alias: "id", Type: "int64", Index: 2, IsBuiltin: true},
			},
		},
	}, q)
}
