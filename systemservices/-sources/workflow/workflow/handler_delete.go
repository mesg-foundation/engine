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
	if err := w.st.Delete(inputs.ID); err != nil {
		return "error", errorOutput{err.Error()}
	}
	return "success", nil
}
