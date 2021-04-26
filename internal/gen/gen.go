package gen

import (
	"bytes"
	"text/template"

	"github.com/nikgalushko/molekula/internal/query"
)

var funcMap = template.FuncMap{
	"sub": func(i int) int {
		return i - 1
	},
	"inc": func(i int) int {
		return i + 1
	},
}

const _map = `
value_{{.Index}}, ok := {{if .IsTop}}data{{else}}raw_value_{{sub .Index}}{{end}}.(map[interface{}]interface{})
if !ok {
	return fmt.Errorf("...")
}

ret_{{.Index}} := make({{.Type}})
for raw_key_{{.Index}}, raw_value_{{.Index}} := range value_{{.Index}} {
	key_{{.Index}}, ok := raw_key_{{.Index}}.({{.KeyType}})
	if !ok {
		return fmt.Errorf("key is not a string")
	}
	{{if .Next.IsBuiltin}}
		value_{{.Index}}, ok := raw_value_{{.Index}}.({{.Next.Type}})
		if !ok {
			return fmt.Errorf("...")
		}

		ret_{{.Index}}[key_{{.Index}}] = value_{{.Index}}
	{{else}}
		{{with .Next}}{{template "T" .}}{{end}}
		ret_{{.Index}}[key_{{.Index}}] = ret_{{inc .Index}}
	{{end}}
}
`

const array = `
value_{{.Index}}, ok := {{if .IsTop}}data{{else}}raw_value_{{sub .Index}}{{end}}.([]interface{})
if !ok {
	return fmt.Errorf("...")
}

ret_{{.Index}} := make({{.Type}}, 0, len(value_{{.Index}}))
for _, raw_value_{{.Index}} := range value_{{.Index}} {
	{{if .Next.IsBuiltin}}
		element, ok := raw_value_{{.Index}}.({{.Next.Type}})
		if !ok {
			return fmt.Errorf("...")
		}
	{{else if .Next.IsArray}}
		{{with .Next}}{{template "TARR" .}}{{end}}
	{{else if .Next.IsMap}}
		{{with .Next}}{{template "TMAP" .}}{{end}}
	{{else if .Next.IsStruct}}
		{{with .Next}}{{template "TSTRUCT" .}}{{end}}
	{{end}}

	ret_{{.Index}} = append(ret_{{.Index}}, {{if .Next.IsBuiltin}}element{{else}}ret_{{inc .Index}}{{end}})
}
`

const builtin = `
ret_0, ok := data.({{.Type}})
if !ok {
	return fmt.Errorf("...")
}
`

const _struct = `
	value_{{.Index}}, ok := {{if .IsTop}}data{{else}}raw_value_{{sub .Index}}{{end}}.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("data")
	}

	ret_{{.Index}} := {{.Type}}{}
	{{range $val := .Fields}}
		{{$val.Name}}_value, ok := value_{{sub .Index}}["{{$val.Alias}}"].({{$val.Type}})
		if !ok {
			return fmt.Errorf("{{$val.Name}}")
		}

		ret_{{sub .Index}}.{{$val.Name}} = {{$val.Name}}_value
	{{end}}
`

const main = `
{{if .IsMap}}
	{{template "TMAP" .}}
{{else if .IsArray}}
	{{template "TARR" .}}
{{else if .IsBuiltin}}
	{{template "TBUILTIN" .}}
{{else if .IsStruct}}
	{{template "TSTRUCT" .}}
{{else}}
	panic("wrong type")
{{end}}
`

// Genrate genrates a function body based on input query.
// It's naive implementation. It's assumed that the parser.Object is valid and fully complies with the specification.
func Generate(q query.Query) (string, error) {
	tmpl := template.Must(template.New("TMAP").Funcs(funcMap).Parse(_map))
	tmpl = template.Must(tmpl.New("TARR").Parse(array))
	tmpl = template.Must(tmpl.New("TBUILTIN").Parse(builtin))
	tmpl = template.Must(tmpl.New("TSTRUCT").Parse(_struct))
	tmpl = template.Must(tmpl.New("T").Parse(main))

	ret := bytes.NewBuffer(nil)

	err := tmpl.Execute(ret, q)
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}
