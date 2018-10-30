package workflow

import mesg "github.com/mesg-foundation/go-service"

// output keys for get task.
const (
	getSuccessOutputKey = "success"
)

// getInputs is the inputs data of get task.
type getInputs struct {
	ID string `json:"id"`
}

// getSuccessOutput is the success output data of get workflow's task.
type getSuccessOutput struct {
	// Workflow details.
	Workflow *WorkflowDocument `json:"workflow"`
}

// getHandler gives the workflow details.
func (w *Workflow) getHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs getInputs
	if err := execution.Data(&inputs); err != nil {
		return newErrorOutput(err)
	}
	workflow, err := w.st.Get(inputs.ID)
	if err != nil {
		return newErrorOutput(err)
	}
	return getSuccessOutputKey, getSuccessOutput{workflow}
}
