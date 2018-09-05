package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var eventDataSchema = map[string]*Parameter{
	"optional": {
		Type:     "String",
		Optional: true,
	},
	"string": {
		Type: "String",
	},
	"number": {
		Type: "Number",
	},
	"boolean": {
		Type: "Boolean",
	},
	"object": {
		Type: "Object",
	},
}

func validateParameterData(parameter string, data interface{}) bool {
	return newParameterValidator("", eventDataSchema[parameter]).Validate(data) == nil
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

func TestValidateParameters(t *testing.T) {
	require.Len(t, validateParametersSchema(eventDataSchema, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}), 0)
	require.Len(t, validateParametersSchema(eventDataSchema, map[string]interface{}{
		"optional": "yeah",
		"string":   "hello",
		"number":   10,
		"boolean":  true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}), 0)
	// 4 errors
	//  - not required string
	//  - invalid number
	//  - invalid boolean
	//  - invalid object
	require.Len(t, validateParametersSchema(eventDataSchema, map[string]interface{}{
		"number":  "string",
		"boolean": 42,
		"object":  false,
	}), 4)
}
