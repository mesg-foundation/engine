package casting

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mesg-foundation/core/protobuf/coreapi"
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

// TaskInputs converts map[string]string to map[string]interface{} based on defined types in the service tasks map.
func TaskInputs(s *coreapi.Service, taskKey string, taskData map[string]string) (map[string]interface{}, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			m := make(map[string]interface{}, len(taskData))
			for key, value := range taskData {
				param, err := findParam(task.Inputs, key)
				if err != nil {
					return nil, err
				}

				newValue, err := taskInputs(value, param.Type)
				if err != nil {
					return nil, fmt.Errorf("task %q - %s", taskKey, err)
				}
				if newValue != nil {
					m[key] = newValue
				}
			}
			return m, nil
		}
	}
	return nil, fmt.Errorf("task %q does not exists", taskKey)
}

// findParam return a param based on the key from a list of parameter
func findParam(parameters []*coreapi.Parameter, key string) (*coreapi.Parameter, error) {
	for _, p := range parameters {
		if p.Key == key {
			return p, nil
		}
	}
	return nil, fmt.Errorf("task input %q does not exists", key)
}

// taskInputs converts single value based on its type.
func taskInputs(value, inputType string) (interface{}, error) {
	c, ok := casters[inputType]
	if !ok {
		return nil, fmt.Errorf("input %q - invalid type", value)
	}
	return c(value)
}
