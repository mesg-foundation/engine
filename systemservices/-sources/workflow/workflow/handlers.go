package workflow

import mesg "github.com/mesg-foundation/go-service"

// WorkflowDocument combines a workflow with additional info.
type WorkflowDocument struct {
	// ID is the unique id for workflow.
	ID string

	// Name is the optionally set unique name for workflow.
	Name string

	// Definition of workflow.
	Definition WorkflowDefinition
}

// output key for errors.
const errOutputKey = "error"

// errorOutput is the error output data.
type errorOutput struct {
	Message string `json:"message"`
}

// newErrorOutput returns a new error output from given err.
func newErrorOutput(err error) (outputKey string, outputData mesg.Data) {
	return errOutputKey, errorOutput{Message: err.Error()}
}
