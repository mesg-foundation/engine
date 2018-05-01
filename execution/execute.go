package execution

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/types"
)

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
