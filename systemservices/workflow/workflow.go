package workflow

import (
	"errors"

	"github.com/mesg-foundation/core/execution"
)

// WSS's tasks.
const (
	CreateTaskKey = "create"
	DeleteTaskKey = "delete"
)

// CreateInputs maps create task's inputs.
// name is optional and has to be unique.
func CreateInputs(yaml []byte, name string) map[string]interface{} {
	return map[string]interface{}{
		"yaml": string(yaml),
		"name": name,
	}
}

// CreateOutputs maps create task's outputs.
// id is unique workflow id.
func CreateOutputs(e *execution.Execution) (id string, err error) {
	switch e.OutputKey {
	case "success":
		return e.OutputData["id"].(string), nil
	case "error":
		return "", errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// DeleteInputs maps delete task's inputs.
// id is unique workflow id.
func DeleteInputs(id string) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

// DeleteOutputs maps delete task's outputs.
func DeleteOutputs(e *execution.Execution) error {
	switch e.OutputKey {
	case "success":
		return nil
	case "error":
		return errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}
