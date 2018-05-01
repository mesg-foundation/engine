package execution

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/types"
)

// Execute moves an exection from the pending to the in progress queue and publish the job for processing
func (execution *Execution) Execute() (reply *types.TaskReply, err error) {
	err = execution.moveFromPendingToInProgress()
	if err != nil {
		return
	}
	execution.ExecutedAt = time.Now()
	log.Println("[PROCESSING]", execution.Task)

	channel := execution.Service.TaskSubscriptionChannel()

	data, err := json.Marshal(execution.Inputs)
	if err != nil {
		return
	}

	reply = &types.TaskReply{
		ExecutionID: execution.ID,
		Task:        execution.Task,
		Data:        string(data),
	}

	go pubsub.Publish(channel, reply)
	return
}
