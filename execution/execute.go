package execution

import (
	"log"
	"time"

	"github.com/mesg-foundation/core/pubsub"
)

// Execute moves an exection from the pending to the in progress queue and publish the job for processing
func (execution *Execution) Execute() (err error) {
	err = execution.moveFromPendingToInProgress()
	if err != nil {
		return
	}
	execution.ExecutedAt = time.Now()
	log.Println("[PROCESSING]", execution.Task)

	channel := execution.Service.TaskSubscriptionChannel()

	go pubsub.Publish(channel, execution)
	return
}
