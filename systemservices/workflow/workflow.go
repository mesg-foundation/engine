package workflow

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/utils/workflowparser"
)

// WSS's tasks.
const (
	CreateTaskKey = "create"
	DeleteTaskKey = "delete"
)

// CreateInputs maps create task's inputs.
// name is optional and has to be unique.
func CreateInputs(definition workflowparser.WorkflowDefinition, name string) (map[string]interface{}, error) {
	// TODO: this hack is not something that we should do but
	// it's needed because *parameterValidator is not able to identify
	// structs for now.
	definitionData, err := json.Marshal(definition)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := json.Unmarshal(definitionData, &data); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"definition": data,
		"name":       name,
	}, nil
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
