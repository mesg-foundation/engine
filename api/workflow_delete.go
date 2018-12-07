package api

import (
	"github.com/mesg-foundation/core/systemservices/workflow"
)

// DeleteWorkflow stops and deletes a workflow. id can be unique id or name.
func (a *API) DeleteWorkflow(id string) (err error) {
	return newWorkflowDeletor(a).Delete(id)
}

// workflowDeletor provides functionalities to stop and delete a workflow.
type workflowDeletor struct {
	api *API
}

// newWorkflowCreator returns a new workflow creator with given api.
func newWorkflowDeletor(api *API) *workflowDeletor {
	return &workflowDeletor{
		api: api,
	}
}

// Delete stops and deletes a workflow. id can be unique id or name.
func (w *workflowDeletor) Delete(id string) (err error) {
	workflowServiceID, err := w.api.systemservices.WorkflowServiceID()
	if err != nil {
		return err
	}
	exec, err := w.api.ExecuteAndListen(workflowServiceID, workflow.DeleteTaskKey, workflow.DeleteInputs(id))
	if err != nil {
		return err
	}
	return workflow.DeleteOutputs(exec)
}
