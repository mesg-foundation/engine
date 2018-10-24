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

	wdoc := &WorkflowDocument{
		ID:         id,
		Name:       inputs.Name,
		Definition: def,
	}

	if err := w.st.Save(wdoc); err != nil {
		return "error", errorOutput{err.Error()}
	}

	if err := w.vm.Run(wdoc); err != nil {
		return "error", errorOutput{err.Error()}
	}

	return "success", createSuccessOutput{id}
}
