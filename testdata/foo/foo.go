package foo

import (
	"fmt"
)

type Foo struct {
	Str    string `molekula:"version"`
	Intrf  interface{}
	ArrInt []int64
	Float  float64
	Int    int64
}

func CallFoo() {
	fmt.Println("foo called")
}
