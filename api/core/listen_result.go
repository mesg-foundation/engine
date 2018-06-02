package core

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/database/services"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenResult will listne for results from a services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	if err = validateTaskKey(&service, request.TaskFilter); err != nil {
		return
	}
	if err = validateOutputKey(&service, request.TaskFilter, request.OutputFilter); err != nil {
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

func validateTaskKey(service *service.Service, taskFilter string) (err error) {
	if taskFilter == "" {
		return
	}
	if taskFilter == "*" {
		return
	}
	_, ok := service.Tasks[taskFilter]
	if ok {
		return
	}
	err = errors.New("Invalid taskFilter: " + taskFilter)
	return
}

func validateOutputKey(service *service.Service, taskKey string, outputFilter string) (err error) {
	if outputFilter == "" || outputFilter == "*" {
		return
	}
	task, ok := service.Tasks[taskKey]
	if !ok {
		err = errors.New("Invalid taskKey: " + taskKey)
		return
	}
	_, ok = task.Outputs[outputFilter]
	if ok {
		return
	}
	err = errors.New("Invalid outputFilter: " + outputFilter)
	return
}

func isSubscribedTask(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Task}, request.TaskFilter)
}

func isSubscribedOutput(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Output}, request.OutputFilter)
}
