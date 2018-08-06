package core

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenResult will listen for results from a services
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) error {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return err
	}
	if err = validateTaskKey(&service, request.TaskFilter); err != nil {
		return err
	}
	if err = validateOutputKey(&service, request.TaskFilter, request.OutputFilter); err != nil {
		return err
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
	return nil
}

func validateTaskKey(service *service.Service, taskKey string) error {
	if taskKey == "" || taskKey == "*" {
		return nil
	}
	_, ok := service.Tasks[taskKey]
	if ok {
		return nil
	}
	return fmt.Errorf("Task %q doesn't exist in this service", taskKey)
}

func validateOutputKey(service *service.Service, taskKey string, outputFilter string) error {
	if outputFilter == "" || outputFilter == "*" {
		return nil
	}
	task, ok := service.Tasks[taskKey]
	if !ok {
		return fmt.Errorf("Task %q doesn't exist in this service", taskKey)
	}
	_, ok = task.Outputs[outputFilter]
	if !ok {
		return fmt.Errorf("Output %q doesn't exist in the task %q of this service", outputFilter, taskKey)
	}
	return nil
}

func isSubscribedTask(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Task}, request.TaskFilter)
}

func isSubscribedOutput(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Output}, request.OutputFilter)
}
