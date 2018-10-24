package workflow

import (
	"errors"

	"github.com/mesg-foundation/core/api"
)

// WSS's tasks.
const (
	createTaskKey = "create"
	deleteTaskKey = "delete"
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

// Create creates and runs a workflow file with an optionally given unique name.
func (w *Workflow) Create(file []byte, name string) (id string, err error) {
	e, err := w.api.ExecuteAndListen(w.serviceID, createTaskKey, map[string]interface{}{
		"file": string(file),
		"name": name,
	})
	if err != nil {
		return "", err
	}

	switch e.Output {
	case "success":
		return e.OutputData["id"].(string), nil
	case "error":
		return "", errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// Delete stops and deletes workflow with id.
func (w *Workflow) Delete(id string) (err error) {
	e, err := w.api.ExecuteAndListen(w.serviceID, deleteTaskKey, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	switch e.Output {
	case "success":
		return nil
	case "error":
		return errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}
