package service

import (
	"testing"

	"github.com/stvp/assert"
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
	assert.Nil(t, event.Data["optional"].Validate("presence"))
	assert.Nil(t, event.Data["optional"].Validate(nil))
	// this parameter is required
	assert.NotNil(t, event.Data["string"].Validate(nil))
}

func TestString(t *testing.T) {
	assert.Nil(t, event.Data["string"].Validate("valid"))
	assert.NotNil(t, event.Data["string"].Validate(false))
}

func TestNumber(t *testing.T) {
	assert.Nil(t, event.Data["number"].Validate(10.5))
	assert.Nil(t, event.Data["number"].Validate(10))
	assert.NotNil(t, event.Data["number"].Validate("not a number"))
}

func TestBoolean(t *testing.T) {
	assert.Nil(t, event.Data["boolean"].Validate(true))
	assert.Nil(t, event.Data["boolean"].Validate(false))
	assert.NotNil(t, event.Data["boolean"].Validate("not a boolean"))
}

func TestObject(t *testing.T) {
	assert.Nil(t, event.Data["object"].Validate(map[string]interface{}{
		"foo": "bar",
	}))
	assert.Nil(t, event.Data["object"].Validate([]interface{}{
		"foo",
		"bar",
	}))
	assert.NotNil(t, event.Data["object"].Validate(42))
}

func TestValidateParameters(t *testing.T) {
	assert.Equal(t, 0, len(validateParameters(event.Data, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	})))
	assert.Equal(t, 0, len(validateParameters(event.Data, map[string]interface{}{
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
	assert.Equal(t, 4, len(validateParameters(event.Data, map[string]interface{}{
		"number":  "string",
		"boolean": 42,
		"object":  false,
	})))
}

func TestValidParameters(t *testing.T) {
	assert.True(t, validParameters(event.Data, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}))
	assert.False(t, validParameters(event.Data, map[string]interface{}{}))
}
