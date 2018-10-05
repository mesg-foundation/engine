package execution

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
)

// Complete marks an execution as complete and puts into the list of processed tasks.
func (execution *Execution) Complete(outputKey string, outputData map[string]interface{}) error {
	var (
		s       = execution.Service
		taskKey = execution.Task
	)
	task, err := s.GetTask(taskKey)
	if err != nil {
		return err
	}
	output, err := task.GetOutput(outputKey)
	if err != nil {
		return err
	}
	if err := output.RequireData(outputData); err != nil {
		return err
	}

	if err := execution.moveFromInProgressToProcessed(); err != nil {
		return err
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = outputKey
	execution.OutputData = outputData

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), execution)
	return nil
}
