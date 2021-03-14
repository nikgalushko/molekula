package main

import (
	"fmt"

	"github.com/nikgalushko/molekula/internal/gen"
	"github.com/nikgalushko/molekula/internal/parser"
)

func main() {
	m, err := parser.Parse("./testdata/foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(m)

	for _, files := range m {
		for _, f := range files {
			for name, fields := range f.Structs {
				code, err := gen.Generate(name, fields)
				if err != nil {
					panic(err)
				}

				fmt.Println(code)
			}
		}
	}
}
