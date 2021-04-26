package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_RawTypeName(t *testing.T) {
	tests := map[string]struct {
		T           Type
		RawTypeName string
	}{
		"simple map": {
			RawTypeName: "map[string]int",
			T: Map{
				Key:   BuiltIn("string"),
				Value: BuiltIn("int"),
			},
		},
		"nested simple map": {
			RawTypeName: "map[string]map[int64]float32",
			T: Map{
				Key: BuiltIn("string"),
				Value: Map{
					Key:   BuiltIn("int64"),
					Value: BuiltIn("float32"),
				},
			},
		},
		"simple map with simple slice as vlaue": {
			RawTypeName: "map[int][]string",
			T: Map{
				Key: BuiltIn("int"),
				Value: Array{
					Element: BuiltIn("string"),
				},
			},
		},
		"slice of simple map": {
			RawTypeName: "[]map[string]int",
			T: Array{
				Element: Map{
					Key:   BuiltIn("string"),
					Value: BuiltIn("int"),
				},
			},
		},
		"slice of nested map with simple slice as value": {
			RawTypeName: "[]map[int]map[string][]float64",
			T: Array{
				Element: Map{
					Key: BuiltIn("int"),
					Value: Map{
						Key: BuiltIn("string"),
						Value: Array{
							Element: BuiltIn("float64"),
						},
					},
				},
			},
		},
	}

	for title, tt := range tests {
		assert.Equal(t, tt.RawTypeName, tt.T.RawTypeName(), title)
	}
}
