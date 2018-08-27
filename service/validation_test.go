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

func validateParemeterData(parameter string, data interface{}) bool {
	return newParameterValidator("", eventDataSchema[parameter]).Validate(data) == nil
}

func TestRequired(t *testing.T) {
	require.True(t, validateParemeterData("optional", "presence"))
	require.True(t, validateParemeterData("optional", nil))
	// this parameter is required
	require.False(t, validateParemeterData("string", nil))
}

func TestString(t *testing.T) {
	require.True(t, validateParemeterData("string", "valid"))
	require.False(t, validateParemeterData("string", false))
}

func TestNumber(t *testing.T) {
	require.True(t, validateParemeterData("number", 10.5))
	require.True(t, validateParemeterData("number", 10))
	require.False(t, validateParemeterData("number", "not a number"))
}

func TestBoolean(t *testing.T) {
	require.True(t, validateParemeterData("boolean", true))
	require.True(t, validateParemeterData("boolean", false))
	require.False(t, validateParemeterData("boolean", "not a boolean"))
}

func TestObject(t *testing.T) {
	require.True(t, validateParemeterData("object", map[string]interface{}{
		"foo": "bar",
	}))
	require.True(t, validateParemeterData("object", []interface{}{
		"foo",
		"bar",
	}))
	require.False(t, validateParemeterData("object", 42))
}

func TestValidateParameters(t *testing.T) {
	s := &Service{}
	require.Len(t, s.ValidateParametersSchema(eventDataSchema, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}), 0)
	require.Len(t, s.ValidateParametersSchema(eventDataSchema, map[string]interface{}{
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
	require.Len(t, s.ValidateParametersSchema(eventDataSchema, map[string]interface{}{
		"number":  "string",
		"boolean": 42,
		"object":  false,
	}), 4)
}
