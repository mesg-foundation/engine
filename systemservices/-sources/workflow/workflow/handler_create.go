package workflow

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/cnf/structhash"
	mesg "github.com/mesg-foundation/go-service"
	uuid "github.com/satori/go.uuid"
)

// output keys for create task.
const (
	createSuccessOutputKey = "success"
)

// createInputs is the input data of new workflow creation.
type createInputs struct {
	File string `json:"file"`
	Name string `json:"name"`
}

// createSuccessOutput is the success output data of new workflow creation.
type createSuccessOutput struct {
	// ID of the workflow.
	ID string `json:"id"`
}

// createHandler creates a new workflow and runs it.
func (w *Workflow) createHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs createInputs
	if err := execution.Data(&inputs); err != nil {
		return newErrorOutput(err)
	}

	def, err := ParseYAML(strings.NewReader(inputs.File))
	if err != nil {
		return newErrorOutput(err)
	}

	id := w.generateHashFromDefinition(def)

	wdoc := &WorkflowDocument{
		ID:         id,
		CreationID: uuid.NewV4().String(),
		Name:       inputs.Name,
		Definition: def,
	}

	if err := w.st.Save(wdoc); err != nil {
		return newErrorOutput(err)
	}

	if err := w.vm.Run(wdoc); err != nil {
		return newErrorOutput(err)
	}

	return createSuccessOutputKey, createSuccessOutput{id}
}

func (w *Workflow) generateHashFromDefinition(def WorkflowDefinition) string {
	h := sha1.New()
	h.Write(structhash.Dump(def, 1))
	return fmt.Sprintf("%x", h.Sum(nil))
}
