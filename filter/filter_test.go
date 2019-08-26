package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		filter Filter
		data   map[string]interface{}
		match  bool
	}{
		{ // not matching filter
			filter: Filter{Conditions: []Condition{
				{Key: "foo", Predicate: EQ, Value: "xx"},
			}},
			data:  map[string]interface{}{"foo": "bar"},
			match: false,
		},
		{ // matching multiple conditions
			filter: Filter{Conditions: []Condition{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "yyy"},
			}},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: true,
		},
		{ // non matching multiple conditions
			filter: Filter{Conditions: []Condition{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "aaa"},
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
