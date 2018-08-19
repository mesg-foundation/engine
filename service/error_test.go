package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type parameterTests []*parameterTest

type parameterTest struct {
	Key   string
	Type  string
	Value interface{}
	Error string
}

func (tests parameterTests) parameterTestsToMapParameter() map[string]*Parameter {
	params := make(map[string]*Parameter)
	for _, test := range tests {
		params[test.Key] = &Parameter{
			Type: test.Type,
		}
	}
	return params
}

func (tests parameterTests) parameterTestsToMapData() map[string]interface{} {
	params := make(map[string]interface{})
	for _, test := range tests {
		params[test.Key] = test.Value
	}
	return params
}

func (tests parameterTests) assert(t *testing.T, err string) {
	for _, test := range tests {
		require.Contains(t, err, "Value of '"+test.Key+"' is "+test.Error)
	}
}

// Test EventNotFoundError
func TestEventNotFoundError(t *testing.T) {
	err := EventNotFoundError{
		Service:  &Service{Name: "TestEventNotFoundError"},
		EventKey: "TestEventNotFoundErrorEventKey",
	}
	require.Equal(t, "Event 'TestEventNotFoundErrorEventKey' not found in service 'TestEventNotFoundError'", err.Error())
}

// Test InvalidEventDataError
func TestInvalidEventDataError(t *testing.T) {
	tests := parameterTests{
		&parameterTest{Key: "keyString", Type: "String", Value: 2323, Error: "not a string"},
		&parameterTest{Key: "keyNumber", Type: "Number", Value: "string", Error: "not a number"},
		&parameterTest{Key: "keyBoolean", Type: "Boolean", Value: "dwdwd", Error: "not a boolean"},
		&parameterTest{Key: "keyObject", Type: "Object", Value: 2323, Error: "not an object or array"},
		&parameterTest{Key: "keyUnknown", Type: "Unknown", Value: "dwdw", Error: "an invalid type"},
		&parameterTest{Key: "keyRequired", Type: "String", Value: nil, Error: "required"},
	}
	err := InvalidEventDataError{
		Event: &Event{
			Data: tests.parameterTestsToMapParameter(),
		},
		Key:  "TestInvalidEventDataErrorEventKey",
		Data: tests.parameterTestsToMapData(),
	}
	require.Contains(t, err.Error(), "Data of event 'TestInvalidEventDataErrorEventKey' is invalid")
	tests.assert(t, err.Error())
}

// Test TaskNotFoundError
func TestTaskNotFoundError(t *testing.T) {
	err := TaskNotFoundError{
		Service: &Service{Name: "TestTaskNotFoundError"},
		TaskKey: "TestTaskNotFoundErrorEventKey",
	}
	require.Equal(t, "Task 'TestTaskNotFoundErrorEventKey' not found in service 'TestTaskNotFoundError'", err.Error())
}

// Test InvalidTaskInputError
func TestInvalidTaskInputError(t *testing.T) {
	tests := parameterTests{
		&parameterTest{Key: "keyString", Type: "String", Value: 2323, Error: "not a string"},
		&parameterTest{Key: "keyNumber", Type: "Number", Value: "string", Error: "not a number"},
		&parameterTest{Key: "keyBoolean", Type: "Boolean", Value: "dwdwd", Error: "not a boolean"},
		&parameterTest{Key: "keyObject", Type: "Object", Value: 2323, Error: "not an object or array"},
		&parameterTest{Key: "keyUnknown", Type: "Unknown", Value: "dwdw", Error: "an invalid type"},
		&parameterTest{Key: "keyRequired", Type: "String", Value: nil, Error: "required"},
	}
	err := InvalidTaskInputError{
		Task: &Task{
			Inputs: tests.parameterTestsToMapParameter(),
		},
		TaskKey: "TestInvalidTaskInputErrorKey",
		Inputs:  tests.parameterTestsToMapData(),
	}
	require.Contains(t, err.Error(), "Inputs of task 'TestInvalidTaskInputErrorKey' are invalid")
	tests.assert(t, err.Error())
}

// Test OutputNotFoundError
func TestOutputNotFoundError(t *testing.T) {
	err := OutputNotFoundError{
		Service:   &Service{Name: "TestOutputNotFoundError"},
		OutputKey: "TestOutputNotFoundErrorEventKey",
	}
	require.Equal(t, "Output 'TestOutputNotFoundErrorEventKey' not found in service 'TestOutputNotFoundError'", err.Error())
}

// Test InvalidOutputDataError
func TestInvalidOutputDataError(t *testing.T) {
	tests := parameterTests{
		&parameterTest{Key: "keyString", Type: "String", Value: 2323, Error: "not a string"},
		&parameterTest{Key: "keyNumber", Type: "Number", Value: "string", Error: "not a number"},
		&parameterTest{Key: "keyBoolean", Type: "Boolean", Value: "dwdwd", Error: "not a boolean"},
		&parameterTest{Key: "keyObject", Type: "Object", Value: 2323, Error: "not an object or array"},
		&parameterTest{Key: "keyUnknown", Type: "Unknown", Value: "dwdw", Error: "an invalid type"},
		&parameterTest{Key: "keyRequired", Type: "String", Value: nil, Error: "required"},
	}
	err := InvalidOutputDataError{
		Output: &Output{
			Data: tests.parameterTestsToMapParameter(),
		},
		Key:  "TestInvalidOutputDataErrorEventKey",
		Data: tests.parameterTestsToMapData(),
	}
	require.Contains(t, err.Error(), "Outputs of task 'TestInvalidOutputDataErrorEventKey' are invalid")
	tests.assert(t, err.Error())
}
