package service

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type caster func(value string) (interface{}, error)

func castString(value string) (interface{}, error) {
	return value, nil
}
func castNumber(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return i, nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("input %q is not a Number type", value)
	}
	return f, nil
}
func castBoolean(value string) (interface{}, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("input %q is not a Boolean type", value)
	}
	return b, nil
}
func castObject(value string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return nil, fmt.Errorf("input %q is not a Object type", value)
	}
	return v, nil
}

var casters = map[string]caster{
	"String":  castString,
	"Number":  castNumber,
	"Boolean": castBoolean,
	"Object":  castObject,
}

// Cast converts map[string]string to map[string]interface{} based on defined types in the service tasks map.
func (s *Service) Cast(taskKey string, taskData map[string]string) (map[string]interface{}, error) {
	task, ok := s.Tasks[taskKey]
	if !ok {
		return nil, &TaskNotFoundError{Service: s, TaskKey: taskKey}
	}

	m := make(map[string]interface{}, len(taskData))
	for key, value := range taskData {
		inputType, ok := task.Inputs[key]
		if !ok {
			return nil, &InputNotFoundError{Service: s, InputKey: key}
		}

		newValue, err := s.cast(value, inputType.Type)
		if err != nil {
			return nil, fmt.Errorf("Task %q - %s", taskKey, err)
		}
		if newValue != nil {
			m[key] = newValue
		}
	}
	return m, nil
}

// cast converts single value based on its type.
func (s *Service) cast(value, inputType string) (interface{}, error) {
	c, ok := casters[inputType]
	if !ok {
		return nil, fmt.Errorf("input %q - invalid type", value)
	}
	return c(value)
}
