package gen

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nikgalushko/molekula/internal/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestGenerate_Builtin(t *testing.T) {
	q := query.Query{
		IsTop:     true,
		IsBuiltin: true,
		Index:     0,
		Type:      "uint",
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "uint"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) (uint, error))(uint(16))
	assert.NoError(t, err)
	assert.Equal(t, uint(16), ret)
}

func TestGenerate_SimpleArray(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsArray: true,
		Type:    "[]int",
		Next:    &query.Query{IsBuiltin: true, Index: 1, Type: "int"},
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "[]int"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) ([]int, error))([]interface{}{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, ret)
}

func TestGenerate_SimpleNestedArray(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsArray: true,
		Type:    "[][]string",
		Next: &query.Query{
			IsArray: true,
			Type:    "[]string",
			Index:   1,
			Next:    &query.Query{Index: 2, IsBuiltin: true, Type: "string"},
		},
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "[][]string"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) ([][]string, error))([]interface{}{
		[]interface{}{"a", "b", "c"},
		[]interface{}{"d", "e", "f"},
	})
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"a", "b", "c"}, {"d", "e", "f"}}, ret)
}

func TestGenerate_SimpleMap(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsMap:   true,
		Type:    "map[int]string",
		KeyType: "int",
		Next:    &query.Query{IsBuiltin: true, Index: 1, Type: "string"},
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "map[int]string"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) (map[int]string, error))(map[interface{}]interface{}{
		1: "one",
		2: "two",
	})
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: "one", 2: "two"}, ret)
}

func TestGenerate_SimpleNestedMap(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsMap:   true,
		Type:    "map[int]map[string]float64",
		KeyType: "int",
		Next: &query.Query{
			Index:   1,
			IsMap:   true,
			Type:    "map[string]float64",
			KeyType: "string",
			Next:    &query.Query{IsBuiltin: true, Index: 2, Type: "float64"},
		},
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "map[int]map[string]float64"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) (map[int]map[string]float64, error))(map[interface{}]interface{}{
		1: map[interface{}]interface{}{"one": 0.1},
		2: map[interface{}]interface{}{"two": 0.2},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[int]map[string]float64{1: {"one": 0.1}, 2: {"two": 0.2}}, ret)
}

func TestGenerate_MapOfArray(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsMap:   true,
		Type:    "map[int][]int64",
		KeyType: "int",
		Next: &query.Query{
			Index:   1,
			IsArray: true,
			Type:    "[]int64",
			Next:    &query.Query{IsBuiltin: true, Index: 2, Type: "int64"},
		},
	}

	s, err := Generate(q)
	assert.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{src: s, typeOfResult: "map[int][]int64"})
	assert.NoError(t, err)

	ret, err := f.(func(interface{}) (map[int][]int64, error))(map[interface{}]interface{}{
		1: []interface{}{int64(0), int64(1)},
		2: []interface{}{int64(1), int64(0)},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[int][]int64{1: {0, 1}, 2: {1, 0}}, ret)
}

func TestGenerate_Struct(t *testing.T) {
	q := query.Query{
		IsTop:    true,
		IsStruct: true,
		Type:     "custom.Foo",
		Fields: []query.Query{
			{Name: "Gender", Alias: "gender", Index: 1, Type: "string"},
			{Name: "ID", Alias: "id", Index: 1, Type: "int64"},
		},
	}

	s, err := Generate(q)
	require.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{
		src:                   s,
		typeOfResult:          "custom.Foo",
		specialTypeDefinition: reflect.ValueOf((*Foo)(nil)),
	})
	require.NoError(t, err)

	ret, err := f.(func(interface{}) (Foo, error))(map[interface{}]interface{}{
		"gender": "m",
		"id":     int64(123),
	})
	require.NoError(t, err)
	assert.Equal(t, Foo{Gender: "m", ID: 123}, ret)
}

func TestGenerate_MapOfStruct(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsMap:   true,
		Type:    "map[string]custom.Foo",
		KeyType: "string",
		Next: &query.Query{
			IsStruct: true,
			Type:     "custom.Foo",
			Index:    1,
			Fields: []query.Query{
				{Name: "Gender", Alias: "gender", Index: 2, Type: "string"},
				{Name: "ID", Alias: "id", Index: 2, Type: "int64"},
			},
		},
	}

	s, err := Generate(q)
	require.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{
		src:                   s,
		typeOfResult:          "map[string]custom.Foo",
		specialTypeDefinition: reflect.ValueOf((*Foo)(nil)),
	})
	require.NoError(t, err)

	ret, err := f.(func(interface{}) (map[string]Foo, error))(map[interface{}]interface{}{
		"first": map[interface{}]interface{}{
			"gender": "m",
			"id":     int64(123),
		},
		"second": map[interface{}]interface{}{
			"gender": "w",
			"id":     int64(456),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]Foo{
		"first":  {Gender: "m", ID: 123},
		"second": {Gender: "w", ID: 456},
	}, ret)
}

func TestGenerate_ArrayOfStruct(t *testing.T) {
	q := query.Query{
		IsTop:   true,
		IsArray: true,
		Type:    "[]custom.Foo",
		KeyType: "string",
		Next: &query.Query{
			IsStruct: true,
			Type:     "custom.Foo",
			Index:    1,
			Fields: []query.Query{
				{Name: "Gender", Alias: "gender", Index: 2, Type: "string"},
				{Name: "ID", Alias: "id", Index: 2, Type: "int64"},
			},
		},
	}

	s, err := Generate(q)
	require.NoError(t, err)

	f, err := buildCallableFunction(buildSettings{
		src:                   s,
		typeOfResult:          "[]custom.Foo",
		specialTypeDefinition: reflect.ValueOf((*Foo)(nil)),
	})
	require.NoError(t, err)

	ret, err := f.(func(interface{}) ([]Foo, error))([]interface{}{
		map[interface{}]interface{}{
			"gender": "m",
			"id":     int64(999),
		},
		map[interface{}]interface{}{
			"gender": "w",
			"id":     int64(888),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, []Foo{
		{Gender: "m", ID: 999},
		{Gender: "w", ID: 888},
	}, ret)
}

type Foo struct {
	Gender string
	ID     int64
}

type buildSettings struct {
	src                   string
	typeOfResult          string
	specialTypeDefinition reflect.Value
}

func buildCallableFunction(s buildSettings) (interface{}, error) {
	pkgTemplate := `
		package foo

		import (
			"fmt"
			"custom"
		)

		var ret %s

		func wrapper(data interface{}) (%s, error) {
			err := pasrse(data)
			return ret, err
		}

		func pasrse(data interface{}) error {
			%s
			ret = ret_0
			return nil
		}
	`
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	custom := make(map[string]map[string]reflect.Value)
	custom["custom"] = make(map[string]reflect.Value)
	custom["custom"]["Foo"] = s.specialTypeDefinition

	i.Use(custom)

	//fmt.Println(fmt.Sprintf(pkgTemplate, s.typeOfResult, s.typeOfResult, s.src))

	_, err := i.Eval(fmt.Sprintf(pkgTemplate, s.typeOfResult, s.typeOfResult, s.src))
	if err != nil {
		return nil, err
	}

	v, err := i.Eval("foo.wrapper")
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}
