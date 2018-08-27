package execution

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// Complete marks an execution as complete and puts into the list of processed tasks.
func (execution *Execution) Complete(output string, data map[string]interface{}) error {
	serviceOutput, outputFound := execution.Service.Tasks[execution.Task].Outputs[output]
	if !outputFound {
		return &service.TaskOutputNotFoundError{
			TaskKey:       execution.Task,
			TaskOutputKey: output,
			ServiceName:   execution.Service.Name,
		}
	}
	warnings := execution.Service.ValidateParametersSchema(serviceOutput.Data, data)
	if len(warnings) > 0 {
		return &service.InvalidTaskOutputError{
			TaskKey:       execution.Task,
			TaskOutputKey: output,
			Warnings:      warnings,
		}
	}
	err := execution.moveFromInProgressToProcessed()
	if err != nil {
		return err
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), execution)
	return nil
}
