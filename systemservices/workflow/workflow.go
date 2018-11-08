package workflow

import (
	"errors"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/execution"
)

// WSS's tasks.
const (
	CreateTaskKey = "create"
	DeleteTaskKey = "delete"
)

// Workflow is a high level wrapper for Workflow System Service.
// It calls the WSS's tasks and reacts to its event through network.
// WSS responsible for managing and running workflows.
type Workflow struct {
	api       *api.API
	serviceID string
}

// New creates a new Workflow for given WSS serviceID and api.
func New(serviceID string, api *api.API) *Workflow {
	return &Workflow{
		api:       api,
		serviceID: serviceID,
	}
}

// CreateInputs maps create task's inputs.
// name is optional and has to be unique.
func CreateInputs(file []byte, name string) map[string]interface{} {
	return map[string]interface{}{
		"file": string(file),
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
