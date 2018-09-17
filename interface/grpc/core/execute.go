package core

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/core"
)

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *core.ExecuteTaskRequest) (*core.ExecuteTaskReply, error) {
	var inputs map[string]interface{}
	if err := json.Unmarshal([]byte(request.InputData), &inputs); err != nil {
		return nil, err
	}

	executionID, err := s.api.ExecuteTask(request.ServiceID, request.TaskKey, inputs, request.ExecutionTags)
	return &core.ExecuteTaskReply{
		ExecutionID: executionID,
	}, err
}
