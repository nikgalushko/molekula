package testdata

import (
	"fmt"
	"go/token"
)

const (
	date = "08.01.1995"

	colour = "green"
)

/*
	molekula:data
*/
type Bar map[string]map[string]int

//molekula:kek
// Foo is not a Bar
type Foo struct {
	Str    string `molekula:"version"`
	Intrf  interface{}
	ArrInt []int64
	Float  float64
	Int    int64
}

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
//
type (
	// An ArrayType node represents an array or slice type.
	ArrayType struct {
		Lbrack token.Pos // position of "["
		Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
		Elt    Expr      // element type
	}

	// A StructType node represents a struct type.
	StructType struct {
		Struct     token.Pos  // position of "struct" keyword
		Fields     *FieldList // list of field declarations
		Incomplete bool       // true if (source) fields are missing in the Fields list
	}

	// Pointer types are represented via StarExpr nodes.

	// A FuncType node represents a function type.
	FuncType struct {
		Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
		Params  *FieldList // (incoming) parameters; non-nil
		Results *FieldList // (outgoing) results; or nil
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  token.Pos  // position of "interface" keyword
		Methods    *FieldList // list of methods
		Incomplete bool       // true if (source) methods are missing in the Methods list
	}

	// A MapType node represents a map type.
	MapType struct {
		Map   token.Pos // position of "map" keyword
		Key   Expr
		Value Expr
	}

	// A ChanType node represents a channel type.
	ChanType struct {
		Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
		Arrow token.Pos // position of "<-" (token.NoPos if there is no "<-")
		Dir   ChanDir   // channel direction
		Value Expr      // value type
	}
)

//molekula:config_version
type Version int

func (v Version) IsProd() bool {
	return v == 5
}

//molekula:weights
type Weights []float64

func (w Weights) Max() float64 {
	ret := w[0]
	for _, v := range w[1:] {
		if v > ret {
			ret = v
		}
	}

	return ret
}

func CallFoo() {
	f := Foo{}
	fmt.Println("foo called", f)
}

type Value struct {
	Gender string
	ID     int64
}

//molekula:config
type Config map[string]Value

//molekula:config2
type Config2 map[string]int

//molekula:slice
type Slice []Value
