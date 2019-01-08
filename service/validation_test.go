package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var eventDataSchema = []*Parameter{
	{
		Key:      "optional",
		Type:     "String",
		Optional: true,
	},
	{
		Key:  "string",
		Type: "String",
	},
	{
		Key:  "number",
		Type: "Number",
	},
	{
		Key:  "boolean",
		Type: "Boolean",
	},
	{
		Key:  "object",
		Type: "Object",
	},
	{
		Key:  "any",
		Type: "Any",
	},
	{
		Key:      "array",
		Type:     "String",
		Repeated: true,
	},
}

func validateParameterData(paramKey string, data interface{}) bool {
	for _, param := range eventDataSchema {
		if param.Key == paramKey {
			return newParameterValidator(param).Validate(data) == nil
		}
	}
	return false
}

func TestRequired(t *testing.T) {
	require.True(t, validateParameterData("optional", "presence"))
	require.True(t, validateParameterData("optional", nil))
	// this parameter is required
	require.False(t, validateParameterData("string", nil))
}

func TestString(t *testing.T) {
	require.True(t, validateParameterData("string", "valid"))
	require.False(t, validateParameterData("string", false))
}

func TestNumber(t *testing.T) {
	require.True(t, validateParameterData("number", 10.5))
	require.True(t, validateParameterData("number", 10))
	require.False(t, validateParameterData("number", "not a number"))
}

func TestBoolean(t *testing.T) {
	require.True(t, validateParameterData("boolean", true))
	require.True(t, validateParameterData("boolean", false))
	require.False(t, validateParameterData("boolean", "not a boolean"))
}

func TestObject(t *testing.T) {
	require.True(t, validateParameterData("object", map[string]interface{}{
		"foo": "bar",
	}))
	require.True(t, validateParameterData("object", []interface{}{
		"foo",
		"bar",
	}))
	require.False(t, validateParameterData("object", 42))
}

func TestAny(t *testing.T) {
	require.True(t, validateParameterData("any", map[string]interface{}{
		"foo": "bar",
	}))
	require.True(t, validateParameterData("any", []interface{}{
		"foo",
		0,
	}))
	require.True(t, validateParameterData("any", 42))
	require.True(t, validateParameterData("any", "string"))
}

func TestArray(t *testing.T) {
	require.True(t, validateParameterData("array", []interface{}{"foo", "bar"}))
	require.True(t, validateParameterData("array", []interface{}{}))
	require.False(t, validateParameterData("array", []interface{}{10}))
	require.False(t, validateParameterData("array", 42))
}

func TestValidateParameters(t *testing.T) {
	tests := []struct {
		data   string
		errors int
	}{
		{
			data: `{
				"string": "hello",
				"number": 10,
				"boolean": true,
				"object": {
					"foo": "bar"
				},
				"any": 0,
				"array": ["foo", "bar"]
			}`,
			errors: 0,
		},
		{
			data: `{
				"optional": "yeah",
				"string": "hello",
				"number": 10,
				"boolean": true,
				"object": {
					"foo": "bar"
				},
				"any": 0,
				"array": ["foo", "bar"]
			}`,
			errors: 0,
		},
		{
			// 5 errors
			//  - not required string
			//  - invalid number
			//  - invalid boolean
			//  - invalid object
			//  - invalid array
			data: `{
				"number": "string",
				"boolean": 42,
				"object": false,
				"any": 0,
				"array": 42
			}`,
			errors: 5,
		},
	}

	for _, test := range tests {
		var data map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(test.data), &data))
		require.Len(t, validateParametersSchema(eventDataSchema, data), test.errors)
	}
}
