package service

import (
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenTask create a stream that will send data for every task to execute
func (s *Server) ListenTask(request *ListenTaskRequest, stream Service_ListenTaskServer) error {
	service, err := services.Get(request.Token)
	if err != nil {
		return err
	}
	subscription := pubsub.Subscribe(service.TaskSubscriptionChannel())
	for data := range subscription {
		execution := data.(*execution.Execution)
		inputs, _ := json.Marshal(execution.Inputs)
		stream.Send(&TaskData{
			ExecutionID: execution.ID,
			TaskKey:     execution.Task,
			InputData:   string(inputs),
		})
	}
	return nil
}
