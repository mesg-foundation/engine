package execution

import (
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// Complete mark an execution as complete and put it in the list of processed tasks
func (execution *Execution) Complete(output string, data map[string]interface{}) error {
	serviceOutput, outputFound := execution.Service.Tasks[execution.Task].Outputs[output]
	if !outputFound {
		return &service.OutputNotFoundError{
			Service:   execution.Service,
			OutputKey: output,
		}
	}
	if !serviceOutput.IsValid(data) {
		return &service.InvalidOutputDataError{
			Output: serviceOutput,
			Key:    output,
			Data:   data,
		}
	}
	err := execution.moveFromInProgressToProcessed()
	if err != nil {
		return err
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data

	log.Println("[COMPLETED]", execution.Task)

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), execution)

	return nil
}
