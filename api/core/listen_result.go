package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenResult will listne for results from a services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	subscription := pubsub.Subscribe(service.ResultSubscriptionChannel())
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
