package core

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"
)

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *ListenResultRequest, stream Core_ListenResultServer) error {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return err
	}
	if err := validateTask(&service, request); err != nil {
		return err
	}

	ctx := stream.Context()
	channel := service.ResultSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-subscription:
			execution := data.(*execution.Execution)
			if isSubscribedToTags(request, execution) && isSubscribedToTask(request, execution) && isSubscribedToOutput(request, execution) {
				outputs, _ := json.Marshal(execution.OutputData)
				if err := stream.Send(&ResultData{
					ExecutionID:   execution.ID,
					TaskKey:       execution.Task,
					OutputKey:     execution.Output,
					OutputData:    string(outputs),
					ExecutionTags: execution.Tags,
				}); err != nil {
					return err
				}
			}
		}
	}
}

func validateTask(service *service.Service, request *ListenResultRequest) error {
	if err := validateTaskKey(service, request.TaskFilter); err != nil {
		return err
	}
	return validateOutputKey(service, request.TaskFilter, request.OutputFilter)
}

func validateTaskKey(service *service.Service, taskKey string) error {
	if taskKey == "" || taskKey == "*" {
		return nil
	}
	if _, ok := service.Tasks[taskKey]; ok {
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

func isSubscribedToTask(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Task}, request.TaskFilter)
}

func isSubscribedToOutput(request *ListenResultRequest, e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Output}, request.OutputFilter)
}

func isSubscribedToTags(request *ListenResultRequest, e *execution.Execution) bool {
	for _, tag := range request.TagFilters {
		if !array.IncludedIn(e.Tags, tag) {
			return false
		}
	}
	return true
}
