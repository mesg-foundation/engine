package workflow

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/cnf/structhash"
	mesg "github.com/mesg-foundation/go-service"
)

// createInputs is the inputs data of workflow creation.
type createInputs struct {
	File string `json:"file"`
	Name string `json:"name"`
}

// createSuccessOutput is the output data of creating a new workflow.
type createSuccessOutput struct {
	// ID of the workflow.
	ID string `json:"id"`
}

// WorkflowDocument combines a workflow with additional info.
type WorkflowDocument struct {
	// ID is the unique id for workflow.
	ID string

	// Name is the optionally set unique name for workflow.
	Name string

	// Definition of workflow.
	Definition WorkflowDefinition
}

// createHandler creates a new workflow and runs it.
func (w *Workflow) createHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs createInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", errorOutput{err.Error()}
	}
	def, err := ParseYAML(strings.NewReader(inputs.File))
	if err != nil {
		return "error", errorOutput{err.Error()}
	}
	h := sha1.New()
	h.Write(structhash.Dump(def, 1))
	id := fmt.Sprintf("%x", h.Sum(nil))
	if err := w.st.Save(&WorkflowDocument{
		ID:         id,
		Name:       inputs.Name,
		Definition: def,
	}); err != nil {
		return "error", errorOutput{err.Error()}
	}
	return "success", createSuccessOutput{id}
}

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

// errorOutput is the error output data.
type errorOutput struct {
	Message string `json:"message"`
}
