package service

import (
	"encoding/json"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// Listen for tasks
func (s *Server) ListenTask(request *ListenTaskRequest, stream Service_ListenTaskServer) (err error) {
	subscription := pubsub.Subscribe(request.Service.TaskSubscriptionChannel())
	for data := range subscription {
		execution := data.(*execution.Execution)
		inputs, _ := json.Marshal(execution.Inputs)
		stream.Send(&TaskData{
			ExecutionID: execution.ID,
			TaskKey:     execution.Task,
			InputData:   string(inputs),
		})
	}
	return
}
