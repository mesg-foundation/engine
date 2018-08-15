package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var serviceTest = &Service{
	Events: map[string]*Event{
		"test": {
			Data: map[string]*Parameter{
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
			},
		},
	},
}
var event = serviceTest.Events["test"]

func TestRequired(t *testing.T) {
	require.Nil(t, event.Data["optional"].Validate("presence"))
	require.Nil(t, event.Data["optional"].Validate(nil))
	// this parameter is required
	require.NotNil(t, event.Data["string"].Validate(nil))
}

func TestString(t *testing.T) {
	require.Nil(t, event.Data["string"].Validate("valid"))
	require.NotNil(t, event.Data["string"].Validate(false))
}

func TestNumber(t *testing.T) {
	require.Nil(t, event.Data["number"].Validate(10.5))
	require.Nil(t, event.Data["number"].Validate(10))
	require.NotNil(t, event.Data["number"].Validate("not a number"))
}

func TestBoolean(t *testing.T) {
	require.Nil(t, event.Data["boolean"].Validate(true))
	require.Nil(t, event.Data["boolean"].Validate(false))
	require.NotNil(t, event.Data["boolean"].Validate("not a boolean"))
}

func TestObject(t *testing.T) {
	require.Nil(t, event.Data["object"].Validate(map[string]interface{}{
		"foo": "bar",
	}))
	require.Nil(t, event.Data["object"].Validate([]interface{}{
		"foo",
		"bar",
	}))
	require.NotNil(t, event.Data["object"].Validate(42))
}

func TestValidateParameters(t *testing.T) {
	require.Equal(t, 0, len(validateParameters(event.Data, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	})))
	require.Equal(t, 0, len(validateParameters(event.Data, map[string]interface{}{
		"optional": "yeah",
		"string":   "hello",
		"number":   10,
		"boolean":  true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	})))
	// 4 errors
	//  - not required string
	//  - invalid number
	//  - invalid boolean
	//  - invalid object
	require.Equal(t, 4, len(validateParameters(event.Data, map[string]interface{}{
		"number":  "string",
		"boolean": 42,
		"object":  false,
	})))
}

func TestValidParameters(t *testing.T) {
	require.True(t, validParameters(event.Data, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}))
	require.False(t, validParameters(event.Data, map[string]interface{}{}))
}
