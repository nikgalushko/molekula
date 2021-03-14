package gen

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
	"text/template"

	"github.com/nikgalushko/molekula/internal/parser"
)

// StructName, StructName, Code
const tmpl = `
	func %sUnmarshal(binMap map[interface{}]interface{}) (ret %s, err error) {
		%s
		return
	}
`

const simpleSnippet = `
		{{.Name}}_raw, ok := binMap["{{.Name}}"]
		if !ok {
			return {{.Struct}}{}, fmt.Errorf("field {{.Name}} doesnt exist")
		}

		{{.Name}}_value, ok := {{.Name}}_raw.({{.Type}})
		if !ok {
			return {{.Struct}}{}, fmt.Errorf("{{.Name}} is not a {{.Type}} is %T", {{.Name}}_raw)
		}

		ret.{{.Field}} = {{.Name}}_value
`

const interfaceSnippet = `
		{{.Name}}_raw, ok := binMap["{{.Name}}"]
		if !ok {
			return {{.Struct}}{}, fmt.Errorf("field {{.Name}} doesnt exist")
		}

		ret.{{.Field}} = {{.Name}}_raw 
`

const simpleArraySnippet = `
		{{.Name}}_raw, ok := binMap["{{.Name}}"]
		if !ok {
			return {{.Struct}}{}, fmt.Errorf("field {{.Name}} doesnt exist")
		}

		{{.Name}}_arr, ok := {{.Name}}_raw.([]{{.Type}})
		if !ok {
			return {{.Struct}}{}, fmt.Errorf("{{.Name}} is not a []{{.Type}} is %T", {{.Name}}_raw)
		}

		ret.{{.Field}} = make([]{{.Type}}, len({{.Name}}_arr))
		for i, elm := range {{.Name}}_arr {
			ret.{{.Field}}[i] = elm
		}
`

func Generate(structName string, fields []parser.Field) (string, error) {
	type Inject struct {
		Name   string
		Struct string
		Type   string
		Field  string
	}
	var code strings.Builder

	for _, f := range fields {
		goType, typeName := typeOf(f.Type)

		var t *template.Template

		switch goType {
		case "simple":
			t = template.Must(template.New("simple-snippet").Parse(simpleSnippet))
		case "interface":
			t = template.Must(template.New("interface-snippet").Parse(interfaceSnippet))
		case "array":
			t = template.Must(template.New("simple-array-snippet").Parse(simpleArraySnippet))
		}

		if t == nil {
			continue
		}

		err := t.Execute(&code, Inject{
			Name:   f.Alias,
			Field:  f.Name,
			Type:   typeName,
			Struct: structName,
		})
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(tmpl, structName, structName, code.String()), nil
}

func typeOf(t ast.Expr) (string, string) {
	fmt.Println("Type:", t, reflect.TypeOf(t))

	simple, ok := t.(*ast.Ident)
	if ok {
		return "simple", simple.Name
	}

	_, ok = t.(*ast.InterfaceType)
	if ok {
		return "interface", "interface{}"
	}

	arr, ok := t.(*ast.ArrayType)
	if ok {
		return "array", arr.Elt.(*ast.Ident).Name
	}

	return "", ""
}
