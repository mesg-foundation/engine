package execution

import (
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
)

// Complete mark an execution as complete and put it in the list of processed tasks
func (execution *Execution) Complete(output string, data interface{}) (err error) {
	err = execution.moveFromInProgressToProcessed()
	if err != nil {
		return
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data

	log.Println("[COMPLETED]", execution.Task)

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), execution)

	return
}
