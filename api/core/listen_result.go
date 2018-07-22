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

// ListenResult will listen for results from a services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) (err error) {
	svc, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	if err = validateTaskKey(&svc, request.TaskFilter); err != nil {
		return
	}
	if err = validateOutputKey(&svc, request.TaskFilter, request.OutputFilter); err != nil {
		return
	}
	subscription := pubsub.Subscribe(svc.ResultSubscriptionChannel())
	for data := range subscription {
		exec := data.(*execution.Execution)
		if isSubscribedTask(request, exec) && isSubscribedOutput(request, exec) {
			outputs, _ := json.Marshal(exec.OutputData)
			stream.Send(&ResultData{
				ExecutionID: exec.ID,
				TaskKey:     exec.Task,
				OutputKey:   exec.Output,
				OutputData:  string(outputs),
			})
		}
	}
	return
}

func validateTaskKey(service *service.Service, taskKey string) (err error) {
	if taskKey == "" || taskKey == "*" {
		return
	}
	_, ok := service.Tasks[taskKey]
	if ok {
		return
	}
	err = errors.New("Task '" + taskKey + "' doesn't exist in this service")
	return
}

func validateOutputKey(service *service.Service, taskKey string, outputFilter string) (err error) {
	if outputFilter == "" || outputFilter == "*" {
		return
	}
	if taskKey == "" {
		err = errors.New("Cannot filter output without specifying a task")
		return
	}
	task, ok := service.Tasks[taskKey]
	if !ok {
		err = errors.New("Task '" + taskKey + "' doesn't exist in this service")
		return
	}
	_, ok = task.Outputs[outputFilter]
	if !ok {
		err = errors.New("Output '" + outputFilter + "' doesn't exist in the task '" + taskKey + "' of this service")
	}
	return
}

func isSubscribedTask(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Task}, request.TaskFilter)
}

func isSubscribedOutput(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Output}, request.OutputFilter)
}
