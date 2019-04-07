package casting

import (
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/stretchr/testify/require"
)

func TestServiceCast(t *testing.T) {
	var tests = []struct {
		service   *definition.Service
		data      map[string]string
		expected  map[string]interface{}
		expectErr bool
	}{
		{
			createTestServcieWithInputs(nil),
			map[string]string{},
			map[string]interface{}{},
			false,
		},
		{
			createTestServcieWithInputs(map[string]string{
				"a": "String",
				"b": "Number",
				"c": "Number",
				"d": "Boolean",
				"e": "repeated String",
				"f": "repeated Number",
				"g": "repeated Number",
				"h": "repeated Boolean",
			}),
			map[string]string{
				"a": "_",
				"b": "1",
				"c": "1.1",
				"d": "true",
				"e": "a,b",
				"f": "1,2",
				"g": "1.1,2.2",
				"h": "false,true",
			},
			map[string]interface{}{
				"a": "_",
				"b": int64(1),
				"c": 1.1,
				"d": true,
				"e": []interface{}{"a", "b"},
				"f": []interface{}{int64(1), int64(2)},
				"g": []interface{}{1.1, 2.2},
				"h": []interface{}{false, true},
			},
			false,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "NoType"}),
			map[string]string{"a": ""},
			map[string]interface{}{},
			true,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "repeated Number"}),
			map[string]string{"a": "0,a"},
			map[string]interface{}{},
			true,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "String"}),
			map[string]string{"b": "_"},
			map[string]interface{}{},
			true,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "Number"}),
			map[string]string{"a": "_"},
			map[string]interface{}{},
			true,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "Boolean"}),
			map[string]string{"a": "_"},
			map[string]interface{}{},
			true,
		},
		{
			createTestServcieWithInputs(map[string]string{"a": "Object"}),
			map[string]string{"a": `{"b":1}`},
			map[string]interface{}{"a": map[string]interface{}{"b": float64(1)}},
			false,
		},
	}

	for _, tt := range tests {
		got, err := TaskInputs(tt.service, "test", tt.data)
		if tt.expectErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Len(t, tt.expected, len(got), "maps len are not equal")
			require.Equal(t, tt.expected, got, "maps are not equal")
		}
	}

	// test if non-existing key returns error
	_, err := TaskInputs(tests[0].service, "_", nil)
	require.Error(t, err)
}

// createTestServcieWithInputs creates test service with given inputs name and type under "test" task key.
func createTestServcieWithInputs(inputs map[string]string) *definition.Service {
	s := &definition.Service{
		Tasks: []*definition.Task{
			{Key: "test"},
		},
	}

	for name, itype := range inputs {
		s.Tasks[0].Inputs = append(s.Tasks[0].Inputs, &definition.Parameter{
			Key:      name,
			Repeated: strings.HasPrefix(itype, "repeated"),
			Type:     strings.TrimPrefix(itype, "repeated "),
		})
	}
	return s
}
