package workflow

import (
	"log"

	mesg "github.com/mesg-foundation/go-service"
	uuid "github.com/satori/go.uuid"
)

type createResponse struct {
	ID string `json:"id"`
}

func (w *Workflow) createHandler(execution *mesg.Execution) (string, mesg.Data) {
	var data interface{}
	if err := execution.Data(&data); err != nil {
		log.Println("err: ", err)
	}
	log.Println("create", data)
	return "success", createResponse{uuid.NewV4().String()}
}

func (w *Workflow) deleteHandler(execution *mesg.Execution) (string, mesg.Data) {
	var data interface{}
	if err := execution.Data(&data); err != nil {
		log.Println("err: ", err)
	}
	log.Println("delete", data)
	return "success", nil
}
