package ast

import "fmt"

// Type is a generic type for all ast types
type Type interface {
	// RawTypeName return a type as string like map[string]int
	RawTypeName() string
}

// Struct is a union of named fields
type Struct struct {
	// Name is a struct name
	Name   string
	Fields []StructField
}

// RawTypeName returns a struct name like Foo or Bar
func (s Struct) RawTypeName() string {
	return s.Name
}

// StructField is a field of struct
type StructField struct {
	// Name is a Gos' name of field
	Name string
	// Alias is a aerospikes' name of fiels
	Alias string
	Type  Type
}

// Map is a mapping a Key to a Value
type Map struct {
	Key   BuiltIn
	Value Type
}

// RawTypeName returns a full type of map like map[int]string
func (m Map) RawTypeName() string {
	return fmt.Sprintf("map[%s]%s", m.Key, m.Value.RawTypeName())
}

// Array is not a array but slice of elements
type Array struct {
	Element Type
}

// RawTypeName returns a full type of array like []float64
func (a Array) RawTypeName() string {
	return fmt.Sprintf("[]%s", a.Element.RawTypeName())
}

// BuiltIn is built-in Go type: int, uint, string, float64, rune etc.
type BuiltIn string

// RawTypeName returns a general type name: int, uint, string etc.
func (b BuiltIn) RawTypeName() string {
	return string(b)
}
