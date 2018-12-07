package api

import (
	"github.com/mesg-foundation/core/systemservices/workflow"
	"github.com/mesg-foundation/core/utils/workflowparser"
)

// CreateWorkflow creates and runs a new workflow with optionally given unique name.
func (a *API) CreateWorkflow(definition workflowparser.WorkflowDefinition, name string) (id string, err error) {
	return newWorkflowCreator(a).Create(definition, name)
}

// workflowCreator provides functionalities to create and run a new workflow.
type workflowCreator struct {
	api *API
}

// newWorkflowCreator returns a new workflow creator with given api.
func newWorkflowCreator(api *API) *workflowCreator {
	return &workflowCreator{
		api: api,
	}
}

// Create creates and runs a new workflow with optionally given unique name.
func (w *workflowCreator) Create(definition workflowparser.WorkflowDefinition, name string) (id string, err error) {
	inputs, err := workflow.CreateInputs(definition, name)
	if err != nil {
		return "", err
	}
	workflowServiceID, err := w.api.systemservices.WorkflowServiceID()
	if err != nil {
		return "", err
	}
	exec, err := w.api.ExecuteAndListen(workflowServiceID, workflow.CreateTaskKey, inputs)
	if err != nil {
		return "", err
	}
	return workflow.CreateOutputs(exec)
}
