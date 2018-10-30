package workflow

import (
	mesg "github.com/mesg-foundation/go-service"
	"github.com/sirupsen/logrus"
)

// output keys for delete task.
const (
	deleteSuccessOutputKey = "success"
)

// deleteInputs is the inputs data of workflow deletion.
type deleteInputs struct {
	ID string `json:"id"`
}

// deleteHandler stops a workflow and deletes it.
func (w *Workflow) deleteHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs deleteInputs
	if err := execution.Data(&inputs); err != nil {
		return newErrorOutput(err)
	}
	workflow, err := w.st.Get(inputs.ID)
	if err != nil {
		return newErrorOutput(err)
	}
	if err := w.vm.Terminate(workflow.ID); err != nil {
		return newErrorOutput(err)
	}
	if err := w.st.Delete(workflow.ID); err != nil {
		return newErrorOutput(err)
	}

	logrus.WithFields(logrus.Fields{
		workflowCreationIDKey: workflow.CreationID,
		workflowLogKey: WorkflowLog{
			Deleted:                       true,
			WorkflowID:                    workflow.ID,
			WorkflowCreationID:            workflow.CreationID,
			WorkflowName:                  workflow.Name,
			WorkflowDefinitionName:        workflow.Definition.Name,
			WorkflowDefinitionDescription: workflow.Definition.Description,
		}}).Info("workflow is deleted")

	return deleteSuccessOutputKey, nil
}
