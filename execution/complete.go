package execution

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// Complete marks an execution as complete and puts into the list of processed tasks.
func (execution *Execution) Complete(outputKey string, outputData map[string]interface{}) error {
	output, ok := execution.Service.Tasks[execution.Task].Outputs[outputKey]
	if !ok {
		return &service.TaskOutputNotFoundError{
			TaskKey:       execution.Task,
			TaskOutputKey: outputKey,
			ServiceName:   execution.Service.Name,
		}
	}
	warnings := execution.Service.ValidateParametersSchema(output.Data, outputData)
	if len(warnings) > 0 {
		return &service.InvalidTaskOutputError{
			TaskKey:       execution.Task,
			TaskOutputKey: outputKey,
			ServiceName:   execution.Service.Name,
			Warnings:      warnings,
		}
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
