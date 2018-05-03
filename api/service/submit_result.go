package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
)

// SubmitResult of an execution
func (s *Server) SubmitResult(context context.Context, request *SubmitResultRequest) (reply *SubmitResultReply, err error) {
	execution := execution.InProgress(request.ExecutionID)
	if execution == nil {
		err = errors.New("No task in progress with the ID " + request.ExecutionID)
		return
	}
	var data interface{}
	err = json.Unmarshal([]byte(request.OutputData), &data)
	if err != nil {
		return
	}
	err = execution.Complete(request.OutputKey, data)
	if err != nil {
		return
	}
	outputData, err := json.Marshal(execution.OutputData)
	if err != nil {
		return
	}
	reply = &SubmitResultReply{
		Error:      "",
		TaskKey:    execution.Task,
		OutputKey:  execution.Output,
		OutputData: string(outputData),
	}
	return
}
