package execution

import (
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
)

// Complete mark an execution as complete and put it in the list of processed tasks
func (execution *Execution) Complete(output string, data map[string]interface{}) (err error) {
	serviceOutput, outputFound := execution.Service.Tasks[execution.Task].Outputs[output]
	if !outputFound {
		return &MissingOutputError{
			Service: execution.Service,
			Output:  output,
		}
	}
	if !serviceOutput.IsValid(data) {
		return &InvalidOutputError{
			Service:  execution.Service,
			Warnings: serviceOutput.Validate(data),
		}
	}
	err = execution.moveFromInProgressToProcessed()
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
