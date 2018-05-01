package execution

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/types"
)

// Complete mark an execution as complete and put it in the list of processed tasks
func (execution *Execution) Complete(output string, data interface{}) (reply *types.ResultReply, err error) {
	err = execution.moveFromInProgressToProcessed()
	if err != nil {
		return
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data

	outputJSON, err := json.Marshal(data)
	if err != nil {
		return
	}

	reply = &types.ResultReply{
		Task:   execution.Task,
		Output: output,
		Data:   string(outputJSON),
	}

	log.Println("[COMPLETED]", execution.Task)

	go pubsub.Publish(execution.Service.ResultSubscriptionChannel(), reply)

	return
}
