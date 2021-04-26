package parser

import (
	goast "go/ast"
	"reflect"
	"strings"

	"github.com/nikgalushko/molekula/internal/ast"
)

type Object struct {
	// BinName is aerospikes' bin name which parsed from tag 'molekula'
	BinName string
	// Type is a type description
	Type ast.Type
}

// Parse returns a list of Objects which tagged 'molekula' in a input file
func Parse(file *goast.File) []Object {
	v := &visitor{}
	goast.Walk(v, file)

	return v.objects
}

type visitor struct {
	objects        []Object
	currentBinName *string
}

func parseBinName(node *goast.GenDecl) (string, bool) {
	if node.Doc == nil || len(node.Doc.List) == 0 {
		return "", false
	}

	for _, v := range node.Doc.List {
		comment := v.Text

		if len(comment) > 2 {
			switch comment[1] {
			case '/':
				comment = comment[2:]
			case '*':
				comment = comment[2 : len(comment)-2]
			}
		}

		for _, comment := range strings.Split(comment, "\n") {
			comment = strings.TrimSpace(comment)

			if strings.HasPrefix(comment, "molekula:") {
				return strings.TrimPrefix(comment, "molekula:"), true
			}
		}
	}

	return "", false
}

func parseStruct(node *goast.StructType) []ast.StructField {
	description := make([]ast.StructField, len(node.Fields.List))

	for i, f := range node.Fields.List {
		description[i] = ast.StructField{
			Type:  pasrseGoASTType(f.Type),
			Name:  f.Names[0].Name,
			Alias: strings.ToLower(f.Names[0].Name),
		}

		if f.Tag != nil {
			tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
			description[i].Alias = tag.Get("molekula")
		}
	}

	return description
}

func pasrseGoASTType(t goast.Expr) ast.Type {
	switch n := t.(type) {
	case *goast.Ident:
		if n.Obj == nil || n.Obj.Decl == nil {
			return ast.BuiltIn(n.Name)
		}

		typeSpec, ok := n.Obj.Decl.(*goast.TypeSpec)
		if !ok {
			return nil
		}

		return ast.Struct{
			Name:   typeSpec.Name.Name,
			Fields: parseStruct(typeSpec.Type.(*goast.StructType)),
		}
	case *goast.InterfaceType:
		return ast.BuiltIn("interface{}")
	case *goast.ArrayType:
		return ast.Array{
			Element: pasrseGoASTType(n.Elt),
		}
	case *goast.MapType:
		return ast.Map{
			Key:   ast.BuiltIn(n.Key.(*goast.Ident).Name),
			Value: pasrseGoASTType(n.Value),
		}
	}
	return nil
}

func (v *visitor) Visit(n goast.Node) goast.Visitor {
	switch node := n.(type) {
	case *goast.GenDecl:
		binName, ok := parseBinName(node)
		if ok {
			v.currentBinName = &binName
		}
	case *goast.TypeSpec:
		if v.currentBinName == nil {
			break
		}
		switch t := node.Type.(type) {
		case *goast.StructType:
			v.objects = append(v.objects, Object{
				BinName: *v.currentBinName,
				Type: ast.Struct{
					Name:   node.Name.Name,
					Fields: parseStruct(t),
				},
			})
			v.currentBinName = nil
		case *goast.MapType, *goast.ArrayType, *goast.Ident:
			v.objects = append(v.objects, Object{
				Type:    pasrseGoASTType(t),
				BinName: *v.currentBinName,
			})
		}
	}

	return v
}
