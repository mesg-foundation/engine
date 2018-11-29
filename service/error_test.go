package service

import (
	"fmt"
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

func (tests parameterTests) parameterTestsToSliceParameters() []*Parameter {
	params := make([]*Parameter, len(tests))
	for i, test := range tests {
		params[i] = &Parameter{
			Key:  test.Key,
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
		require.Contains(t, err, fmt.Sprintf("Value of %q is %s", test.Key, test.Error))
	}
}

func newParameterTestCases() parameterTests {
	return parameterTests{
		&parameterTest{Key: "keyString", Type: "String", Value: 2323, Error: "not a string"},
		&parameterTest{Key: "keyNumber", Type: "Number", Value: "string", Error: "not a number"},
		&parameterTest{Key: "keyBoolean", Type: "Boolean", Value: "dwdwd", Error: "not a boolean"},
		&parameterTest{Key: "keyObject", Type: "Object", Value: 2323, Error: "not an object or array"},
		&parameterTest{Key: "keyUnknown", Type: "Unknown", Value: "dwdw", Error: "an invalid type"},
		&parameterTest{Key: "keyRequired", Type: "String", Value: nil, Error: "required"},
	}
}

// Test EventNotFoundError
func TestEventNotFoundError(t *testing.T) {
	err := EventNotFoundError{
		EventKey:    "TestEventNotFoundErrorEventKey",
		ServiceName: "TestEventNotFoundError",
	}
	require.Equal(t, `Event "TestEventNotFoundErrorEventKey" not found in service "TestEventNotFoundError"`,
		err.Error())
}

// Test InvalidEventDataError
func TestInvalidEventDataError(t *testing.T) {
	tests := newParameterTestCases()
	err := InvalidEventDataError{
		EventKey:    "TestInvalidEventDataErrorEventKey",
		ServiceName: "TestInvalidEventDataError",
		Warnings: validateParametersSchema(tests.parameterTestsToSliceParameters(),
			tests.parameterTestsToMapData()),
	}
	require.Contains(t, err.Error(), `Data of event "TestInvalidEventDataErrorEventKey" is invalid in service "TestInvalidEventDataError"`)
	tests.assert(t, err.Error())
}

// Test TaskNotFoundError
func TestTaskNotFoundError(t *testing.T) {
	err := TaskNotFoundError{
		TaskKey:     "TestTaskNotFoundErrorEventKey",
		ServiceName: "TestTaskNotFoundError",
	}
	require.Equal(t, `Task "TestTaskNotFoundErrorEventKey" not found in service "TestTaskNotFoundError"`, err.Error())
}

// Test InvalidTaskInputError
func TestInvalidTaskInputError(t *testing.T) {
	tests := newParameterTestCases()
	err := InvalidTaskInputError{
		TaskKey:     "TestInvalidTaskInputErrorKey",
		ServiceName: "TestInvalidTaskInputError",
		Warnings: validateParametersSchema(tests.parameterTestsToSliceParameters(),
			tests.parameterTestsToMapData()),
	}
	require.Contains(t, err.Error(), `Inputs of task "TestInvalidTaskInputErrorKey" are invalid in service "TestInvalidTaskInputError"`)
	tests.assert(t, err.Error())
}

// Test OutputNotFoundError
func TestOutputNotFoundError(t *testing.T) {
	err := TaskOutputNotFoundError{
		TaskKey:       "TaskKey",
		TaskOutputKey: "OutputKey",
		ServiceName:   "TestOutputNotFoundError",
	}
	require.Equal(t, `Output "OutputKey" of task "TaskKey" not found in service "TestOutputNotFoundError"`, err.Error())
}

// Test InvalidOutputDataError
func TestInvalidOutputDataError(t *testing.T) {
	tests := newParameterTestCases()
	err := InvalidTaskOutputError{
		TaskKey:       "TaskKey",
		TaskOutputKey: "OutputKey",
		ServiceName:   "TestInvalidOutputDataError",
		Warnings: validateParametersSchema(tests.parameterTestsToSliceParameters(),
			tests.parameterTestsToMapData()),
	}
	require.Contains(t, err.Error(), `Outputs "OutputKey" of task "TaskKey" are invalid in service "TestInvalidOutputDataError"`)
	tests.assert(t, err.Error())
}

// Test InputNotFoundError
func TestInputNotFoundError(t *testing.T) {
	err := TaskInputNotFoundError{
		TaskKey:      "TaskKey",
		TaskInputKey: "InputKey",
		ServiceName:  "InputNotFoundError",
	}
	require.Equal(t, `Input "InputKey" of task "TaskKey" not found in service "InputNotFoundError"`, err.Error())
}
