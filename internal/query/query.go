package query

import (
	"github.com/nikgalushko/molekula/internal/ast"
	"github.com/nikgalushko/molekula/internal/parser"
)

// Query is start point to code generator
type Query struct {
	// IsTop
	IsTop bool
	// Index is a query nesting level
	Index     int
	IsArray   bool
	IsMap     bool
	IsBuiltin bool
	IsStruct  bool
	// Fields is not empty is IsStruct is true
	Fields []Query
	// Name is name of struct field
	Name string
	// Alias is an alias of struct fields
	Alias string
	// Type is result of call .RawTypeName() function
	Type string
	// KeyType is not empty if IsMap is true
	KeyType string
	// Next is pointer to description of nested type
	Next *Query
}

// Build builds Query from parser.Object for generator.
// It's naive implementation. It's assumed that the parser.Object is valid and fully complies with the specification.
func Build(o parser.Object) Query {
	root := Query{IsTop: true}
	q := &root
	t := o.Type
	index := 0

	for {
		q.Index = index
		q.Type = t.RawTypeName()

		switch kind := t.(type) {
		case ast.BuiltIn:
			q.IsBuiltin = true
			t = nil
		case ast.Array:
			q.IsArray = true
			t = kind.Element
		case ast.Map:
			q.IsMap = true
			q.KeyType = kind.Key.RawTypeName()
			t = kind.Value
		case ast.Struct:
			q.IsStruct = true
			for _, f := range kind.Fields {
				q.Fields = append(q.Fields, Query{
					Name:      f.Name,
					Alias:     f.Alias,
					IsBuiltin: true,
					Index:     index + 1,
					Type:      f.Type.RawTypeName(),
				})
			}
			t = nil
		}

		if t == nil {
			break
		}

		next := &Query{}

		q.Next = next

		q = next
		index++
	}

	return root
}
