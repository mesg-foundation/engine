package core

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/database/services"
	service "github.com/mesg-foundation/core/service"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenResult will listne for results from a services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	if err = validateTaskKey(&service, request.TaskKey); err != nil {
		return
	}
	if err = validateOutputKey(&service, request.TaskKey, request.OutputKey); err != nil {
		return
	}
	subscription := pubsub.Subscribe(service.ResultSubscriptionChannel())
	for data := range subscription {
		execution := data.(*execution.Execution)
		if isSubscribedTask(request, execution) && isSubscribedOutput(request, execution) {
			outputs, _ := json.Marshal(execution.OutputData)
			stream.Send(&ResultData{
				ExecutionID: execution.ID,
				TaskKey:     execution.Task,
				OutputKey:   execution.Output,
				OutputData:  string(outputs),
			})
		}
	}
	return
}

func validateTaskKey(service *service.Service, taskKey string) (err error) {
	if taskKey == "" {
		return
	}
	if taskKey == "*" {
		return
	}
	_, ok := service.Tasks[taskKey]
	if ok {
		return
	}
	err = errors.New("Invalid taskKey: " + taskKey)
	return
}

func validateOutputKey(service *service.Service, taskKey string, outputKey string) (err error) {
	if outputKey == "" {
		return
	}
	if outputKey == "*" {
		return
	}
	task, ok := service.Tasks[taskKey]
	if !ok {
		err = errors.New("Invalid taskKey: " + taskKey)
		return
	}
	_, ok = task.Outputs[outputKey]
	if ok {
		return
	}
	err = errors.New("Invalid outputKey: " + outputKey)
	return
}

func isSubscribedTask(request *ListenResultRequest, e *execution.Execution) bool {
	return includedIn([]string{"", "*", e.Task}, request.TaskKey)
}

func isSubscribedOutput(request *ListenResultRequest, e *execution.Execution) bool {
	return includedIn([]string{"", "*", e.Output}, request.OutputKey)
}

func includedIn(arr []string, value string) bool {
	if len(arr) == 0 {
		return false
	}
	i := 0
	for _, item := range arr {
		if item == value {
			break
		}
		i++
	}
	return i != len(arr)
}
