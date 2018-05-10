package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// Listen for results from the services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) (err error) {
	subscription := pubsub.Subscribe(request.Service.ResultSubscriptionChannel())
	for data := range subscription {
		execution := data.(*execution.Execution)
		outputs, _ := json.Marshal(execution.OutputData)
		stream.Send(&ResultData{
			ExecutionID: execution.ID,
			TaskKey:     execution.Task,
			OutputKey:   execution.Output,
			OutputData:  string(outputs),
		})
	}
	return
}
