package event

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serviceTest = &service.Service{
	Events: map[string]*service.Event{
		"test": &service.Event{
			Data: map[string]*service.Parameter{
				"optional": &service.Parameter{
					Type:     "String",
					Optional: true,
				},
				"string": &service.Parameter{
					Type: "String",
				},
				"number": &service.Parameter{
					Type: "Number",
				},
				"boolean": &service.Parameter{
					Type: "Boolean",
				},
				"object": &service.Parameter{
					Type: "Object",
				},
			},
		},
	},
}
var parameters = serviceTest.Events["test"].Data

func TestEventExists(t *testing.T) {
	assert.True(t, exists(serviceTest, "test"))
	assert.False(t, exists(serviceTest, "test-that-doesnt-exists"))
}

func TestRequired(t *testing.T) {
	assert.Nil(t, checkParameterWarning("optional", parameters["optional"], map[string]interface{}{
		"optional": "presence",
	}))
	assert.Nil(t, checkParameterWarning("optional", parameters["optional"], map[string]interface{}{}))
	// this parameter is required
	assert.NotNil(t, checkParameterWarning("string", parameters["string"], map[string]interface{}{}))
}

func TestString(t *testing.T) {
	assert.Nil(t, checkParameterWarning("string", parameters["string"], map[string]interface{}{
		"string": "valid",
	}))
	assert.NotNil(t, checkParameterWarning("string", parameters["string"], map[string]interface{}{
		"string": false,
	}))
}

func TestNumber(t *testing.T) {
	assert.Nil(t, checkParameterWarning("number", parameters["number"], map[string]interface{}{
		"number": 10.5,
	}))
	assert.Nil(t, checkParameterWarning("number", parameters["number"], map[string]interface{}{
		"number": 10,
	}))
	assert.NotNil(t, checkParameterWarning("number", parameters["number"], map[string]interface{}{
		"number": "not a number",
	}))
}

func TestBoolean(t *testing.T) {
	assert.Nil(t, checkParameterWarning("boolean", parameters["boolean"], map[string]interface{}{
		"boolean": true,
	}))
	assert.Nil(t, checkParameterWarning("boolean", parameters["boolean"], map[string]interface{}{
		"boolean": false,
	}))
	assert.NotNil(t, checkParameterWarning("boolean", parameters["boolean"], map[string]interface{}{
		"boolean": "not a boolean",
	}))
}

func TestObject(t *testing.T) {
	assert.Nil(t, checkParameterWarning("object", parameters["object"], map[string]interface{}{
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}))
	assert.Nil(t, checkParameterWarning("object", parameters["object"], map[string]interface{}{
		"object": []interface{}{
			"foo",
			"bar",
		},
	}))
	assert.NotNil(t, checkParameterWarning("object", parameters["object"], map[string]interface{}{
		"object": 42,
	}))
}

func TestParametersWarnings(t *testing.T) {
	assert.Equal(t, 0, len(parametersWarnings(parameters, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	})))
	assert.Equal(t, 0, len(parametersWarnings(parameters, map[string]interface{}{
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
	assert.Equal(t, 4, len(parametersWarnings(parameters, map[string]interface{}{
		"number":  "string",
		"boolean": 42,
		"object":  false,
	})))
}

func TestValidParameters(t *testing.T) {
	assert.True(t, validParameters(parameters, map[string]interface{}{
		"string":  "hello",
		"number":  10,
		"boolean": true,
		"object": map[string]interface{}{
			"foo": "bar",
		},
	}))
	assert.False(t, validParameters(parameters, map[string]interface{}{}))
}
