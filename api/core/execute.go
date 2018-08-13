package core

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
)

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (*ExecuteTaskReply, error) {
	srv, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	inputs, err := getData(request)
	if err != nil {
		return nil, err
	}
	if err := checkServiceStatus(&srv); err != nil {
		return nil, err
	}
	executionID, err := execute(&srv, request.TaskKey, inputs)
	return &ExecuteTaskReply{
		ExecutionID: executionID,
	}, err
}

func checkServiceStatus(srv *service.Service) error {
	status, err := srv.Status()
	if err != nil {
		return err
	}
	if status != service.RUNNING {
		return &NotRunningServiceError{ServiceID: srv.Hash()}
	}
	return nil
}

func getData(request *ExecuteTaskRequest) (map[string]interface{}, error) {
	var inputs map[string]interface{}
	err := json.Unmarshal([]byte(request.InputData), &inputs)
	return inputs, err
}

func execute(srv *service.Service, key string, inputs map[string]interface{}) (executionID string, err error) {
	exc, err := execution.Create(srv, key, inputs)
	if err != nil {
		return "", err
	}
	return exc.ID, exc.Execute()
}
