package process

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		filter Process_Node_Filter
		data   map[string]interface{}
		match  bool
	}{
		{ // not matching filter
			filter: Process_Node_Filter{Conditions: []Process_Node_Filter_Condition{
				{Key: "foo", Predicate: Process_Node_Filter_Condition_EQ, Value: "xx"},
			}},
			data:  map[string]interface{}{"foo": "bar"},
			match: false,
		},
		{ // matching multiple conditions
			filter: Process_Node_Filter{Conditions: []Process_Node_Filter_Condition{
				{Key: "foo", Predicate: Process_Node_Filter_Condition_EQ, Value: "bar"},
				{Key: "xxx", Predicate: Process_Node_Filter_Condition_EQ, Value: "yyy"},
			}},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: true,
		},
		{ // non matching multiple conditions
			filter: Process_Node_Filter{Conditions: []Process_Node_Filter_Condition{
				{Key: "foo", Predicate: Process_Node_Filter_Condition_EQ, Value: "bar"},
				{Key: "xxx", Predicate: Process_Node_Filter_Condition_EQ, Value: "aaa"},
			}},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: false,
		},
	}
	for i, test := range tests {
		match := test.filter.Match(test.data)
		require.Equal(t, test.match, match, i)
	}
}
