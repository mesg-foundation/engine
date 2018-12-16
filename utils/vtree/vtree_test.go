package vtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Simple struct {
	A string
	B bool
}

var str *string

func TestAnalyse(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		result Value
	}{
		{
			"nil",
			nil,
			Value{},
		},
		{
			"string",
			str,
			Value{Type: Nil},
		},
		{
			"string nil",
			"1",
			Value{Type: String},
		},
		{
			"string slice",
			[]string{"1", "2"},
			Value{Type: Array, Values: []Value{{Type: String}, {Type: String}}},
		},
		{
			"slice of string slices",
			[][]string{{"1"}, {"1", "2"}},
			Value{Type: Array, Values: []Value{
				{Type: Array, Values: []Value{{Type: String}}},
				{Type: Array, Values: []Value{{Type: String}, {Type: String}}},
			}},
		},
		{
			"interface slice",
			[]interface{}{"1", 2, true, nil},
			Value{Type: Array, Values: []Value{
				{Type: String},
				{Type: Number},
				{Type: Bool},
				{Type: Nil},
			}},
		},
		{
			"map",
			map[string]interface{}{
				"a": 1,
				"b": true,
				"c": &Simple{"1", true},
				"d": []*Simple{{"1", true}, {"1", true}},
			},
			Value{Type: Object, Values: []Value{
				{Key: "a", Type: Number},
				{Key: "b", Type: Bool},
				{Key: "c", Type: Object, Values: []Value{
					{Key: "A", Type: String},
					{Key: "B", Type: Bool},
				}},
				{Key: "d", Type: Array, Values: []Value{
					{Type: Object, Values: []Value{
						{Key: "A", Type: String},
						{Key: "B", Type: Bool},
					}},
					{Type: Object, Values: []Value{
						{Key: "A", Type: String},
						{Key: "B", Type: Bool},
					}},
				}},
			}},
		},
		{
			"struct",
			Simple{"1", true},
			Value{Type: Object, Values: []Value{{Key: "A", Type: String}, {Key: "B", Type: Bool}}},
		},
	}

	for _, test := range tests {
		v := Analyze(test.value)
		require.Equal(t, test.result, v, test.name)
	}
}
