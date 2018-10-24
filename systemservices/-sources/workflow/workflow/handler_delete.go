package workflow

import mesg "github.com/mesg-foundation/go-service"

// deleteInputs is the inputs data of workflow deletion.
type deleteInputs struct {
	ID string `json:"id"`
}

// deleteHandler stops a workflow and deletes it.
func (w *Workflow) deleteHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs deleteInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", errorOutput{err.Error()}
	}
	workflow, err := w.st.Get(inputs.ID)
	if err != nil {
		return "error", errorOutput{err.Error()}
	}
	if err := w.vm.Terminate(workflow.ID); err != nil {
		return "error", errorOutput{err.Error()}
	}
	if err := w.st.Delete(workflow.ID); err != nil {
		return "error", errorOutput{err.Error()}
	}
	return "success", nil
}
