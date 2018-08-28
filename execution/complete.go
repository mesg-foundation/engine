package execution

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
)

// Complete marks an execution as complete and puts into the list of processed tasks.
func (execution *Execution) Complete(output string, data map[string]interface{}) error {
	if err := execution.moveFromInProgressToProcessed(); err != nil {
		return err
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), execution)
	return nil
}
