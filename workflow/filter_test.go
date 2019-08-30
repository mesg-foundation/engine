package workflow

// func TestMatch(t *testing.T) {
// 	var tests = []struct {
// 		filters TriggerFilters
// 		data    map[string]interface{}
// 		match   bool
// 	}{
// 		{ // not matching filter
// 			filters: []*TriggerFilter{
// 				{Key: "foo", Predicate: EQ, Value: "xx"},
// 			},
// 			data:  map[string]interface{}{"foo": "bar"},
// 			match: false,
// 		},
// 		{ // matching multiple filters
// 			filters: []*TriggerFilter{
// 				{Key: "foo", Predicate: EQ, Value: "bar"},
// 				{Key: "xxx", Predicate: EQ, Value: "yyy"},
// 			},
// 			data: map[string]interface{}{
// 				"foo": "bar",
// 				"xxx": "yyy",
// 				"aaa": "bbb",
// 			},
// 			match: true,
// 		},
// 		{ // non matching multiple filters
// 			filters: []*TriggerFilter{
// 				{Key: "foo", Predicate: EQ, Value: "bar"},
// 				{Key: "xxx", Predicate: EQ, Value: "aaa"},
// 			},
// 			data: map[string]interface{}{
// 				"foo": "bar",
// 				"xxx": "yyy",
// 				"aaa": "bbb",
// 			},
// 			match: false,
// 		},
// 	}
// 	for i, test := range tests {
// 		match := test.filters.Match(test.data)
// 		assert.Equal(t, test.match, match, i)
// 	}
// }
