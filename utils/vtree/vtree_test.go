package vtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyse(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		result Value
	}{
		{
			"nil",
			nil,
			Value{Type: Nil},
		},
		{
			"string",
			"1",
			Value{Type: String},
		},
		{
			"string slice",
			[]string{"1", "2"},
			Value{Type: Array, Values: []Value{{Type: String}, {Type: String}}},
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
				"c": map[string]interface{}{"A": "1", "B": true},
				"d": []map[string]interface{}{{"A": "1", "B": true}, {"A": "2", "B": true}},
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
	}

	for _, test := range tests {
		v := Analyze(test.value)
		require.Equal(t, test.result, v, test.name)
	}
}

func TestGetByKey(t *testing.T) {
	v := Analyze(map[string]interface{}{
		"A": "1",
		"B": map[string]interface{}{
			"C": "2",
		},
	})

	vv, ok := v.GetByKey("a", true)
	require.False(t, ok)
	require.Equal(t, Value{}, vv)

	vv, ok = v.GetByKey("a", false)
	require.True(t, ok)
	require.Equal(t, Value{Type: String, Key: "A"}, vv)

	vv, ok = v.GetByKey("B", true)
	require.True(t, ok)
	require.Equal(t, Value{Type: Object, Key: "B", Values: []Value{
		{Type: String, Key: "C"},
	}}, vv)
}
