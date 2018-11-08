package workflow

import (
	"log"

	mesg "github.com/mesg-foundation/go-service"
	uuid "github.com/satori/go.uuid"
)

// createResponse is the response message of creating a new workflow.
type createResponse struct {
	// ID of the workflow.
	ID string `json:"id"`
}

// createHandler creates a new workflow and runs it.
func (w *Workflow) createHandler(execution *mesg.Execution) (string, mesg.Data) {
	var data interface{}
	if err := execution.Data(&data); err != nil {
		log.Println("err: ", err)
	}
	log.Println("create", data)
	return "success", createResponse{uuid.NewV4().String()}
}

// deleteHandler stops a workflow and deletes it.
func (w *Workflow) deleteHandler(execution *mesg.Execution) (string, mesg.Data) {
	var data interface{}
	if err := execution.Data(&data); err != nil {
		log.Println("err: ", err)
	}
	log.Println("delete", data)
	return "success", nil
}
