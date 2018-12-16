package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		params   []*Parameter
		data     interface{}
		warnings []*ParameterWarning
		err      error
	}{
		{
			"parameters with valid map data",
			[]*Parameter{
				{Key: "a", Type: "String"},
				{Key: "b", Type: "Number"},
				{Key: "c", Type: "Boolean"},
				{Key: "d", Type: "Any"},
				{Key: "e", Parameters: []*Parameter{
					{Key: "f", Type: "String"},
				}},
				{Key: "g", Optional: true, Type: "Boolean"},
			},
			map[string]interface{}{
				"a": "1",
				"b": 2,
				"c": true,
				"d": true,
				"e": map[string]interface{}{"f": "3"},
			},
			nil,
			nil,
		},
		{
			"parameters with valid struct data",
			[]*Parameter{
				{Key: "a", Type: "String"},
				{Key: "b", Type: "Number"},
				{Key: "c", Type: "Number"},
				{Key: "d", Type: "Boolean"},
				{Key: "e", Type: "Any"},
				{Key: "f", Parameters: []*Parameter{
					{Key: "g", Type: "String"},
				}},
				{Key: "h", Parameters: []*Parameter{
					{Key: "g", Type: "String"},
				}},
			},
			struct {
				A string
				B int
				C float64
				D bool
				E interface{}
				F struct{ G string }
				H interface{}
			}{"1", 2, 0.3, true, 4, struct{ G string }{"6"}, struct{ G string }{"7"}},
			nil,
			nil,
		},
		{
			"repeated parameters with valid data",
			[]*Parameter{
				{Key: "a", Repeated: true, Type: "String"},
				{Key: "b", Repeated: true, Type: "String"},
				{Key: "c", Repeated: true, Parameters: []*Parameter{
					{Key: "d", Type: "String"},
				}},
				{Key: "e", Parameters: []*Parameter{
					{Key: "f", Repeated: true, Parameters: []*Parameter{
						{Key: "g", Type: "String"},
					}},
				}},
			},
			map[string]interface{}{
				"a": []string{"1", "2"},
				"b": []interface{}{"3", "4"},
				"c": []map[string]interface{}{
					{"d": "5"},
					{"d": "6"},
				},
				"e": map[string]interface{}{
					"f": []interface{}{},
				},
			},
			nil,
			nil,
		},
		{
			"parameters with invalid map data",
			[]*Parameter{
				{Key: "a", Type: "String"},
				{Key: "b", Type: "Number"},
				{Key: "c", Type: "Boolean"},
				{Key: "d", Type: "Any"},
				{Key: "e", Parameters: []*Parameter{
					{Key: "f", Type: "String"},
				}},
				{Key: "g", Type: "Invalid"},
				{Key: "h", Type: "String"},
				{Key: "j", Type: "String"},
			},
			map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
				"d": true,
				"e": map[string]interface{}{"f": 3},
				"g": "1",
				"j": []string{"1", "2"},
			},
			[]*ParameterWarning{
				notAStringWarning("a"),
				notANumberWarning("b"),
				notAStringWarning("f"),
				unKnownTypeWarning("g"),
				requiredWarning("h"),
				notAStringWarning("j"),
			},
			nil,
		},
		{
			"repeated parameters with invalid data",
			[]*Parameter{
				{Key: "a", Repeated: true, Type: "String"},
				{Key: "b", Repeated: true, Parameters: []*Parameter{
					{Key: "c", Repeated: true, Parameters: []*Parameter{
						{Key: "d", Type: "Boolean"},
					}},
				}},
				{Key: "e", Parameters: []*Parameter{
					{Key: "f", Type: "Boolean"},
				}},
				{Key: "g", Type: "Boolean"},
				{Key: "h", Parameters: []*Parameter{
					{Key: "j", Type: "Boolean"},
				}},
			},
			map[string]interface{}{
				"a": "1",
				"b": []map[string]interface{}{
					{"c": true},
				},
				"e": nil,
				"g": "2",
				"h": "3",
			},
			[]*ParameterWarning{
				notAnArrayWarning("a"),
				notAnArrayWarning("c"),
				requiredWarning("e"),
				notABooleanWarning("g"),
				notAnObjectWarning("h"),
			},
			nil,
		},
		{
			"complex parameter data",
			[]*Parameter{{
				Key: "article",
				Parameters: []*Parameter{
					{Key: "id", Type: "String"},
					{Key: "location", Parameters: []*Parameter{
						{Key: "city", Type: "String"},
					}},
					{Key: "createdAt", Type: "String"},
				},
			}},
			map[string]interface{}{
				"article": map[string]interface{}{
					"id": bson.NewObjectId(),
					"location": map[string]interface{}{
						"city": "Ankara",
					},
					"createdAt": time.Now(),
				},
			},
			nil,
			nil,
		},
	}

	for _, test := range tests {
		warnings, err := newParameterValidator().Validate(test.params, test.data)
		require.Equal(t, test.err, err, test.name)
		require.Equal(t, len(test.warnings), len(warnings), test.name)
		for i, w := range warnings {
			require.Equal(t, test.warnings[i].Key, w.Key, test.name)
			require.Equal(t, test.warnings[i].Warning, w.Warning, test.name)
		}
	}
}

func notAStringWarning(key string) *ParameterWarning {
	return warning(key, "not a string")
}

func notANumberWarning(key string) *ParameterWarning {
	return warning(key, "not a number")
}

func notABooleanWarning(key string) *ParameterWarning {
	return warning(key, "not a boolean")
}

func notAnArrayWarning(key string) *ParameterWarning {
	return warning(key, "not an array")
}

func notAnObjectWarning(key string) *ParameterWarning {
	return warning(key, "not an object")
}

func unKnownTypeWarning(key string) *ParameterWarning {
	return warning(key, "unknown type")
}

func requiredWarning(key string) *ParameterWarning {
	return warning(key, "required")
}

func warning(key, warning string) *ParameterWarning {
	return &ParameterWarning{Key: key, Warning: warning}
}
