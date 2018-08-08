package execution

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/mesg-foundation/core/pubsub"
)

// Execute moves an execution from pending to in progress queue and publishes the job for processing.
func (execution *Execution) Execute() error {
	if err := execution.moveFromPendingToInProgress(); err != nil {
		return err
	}
	execution.ExecutedAt = time.Now()
	log.Info("[PROCESSING]", execution.Task)

	channel := execution.Service.TaskSubscriptionChannel()

	go pubsub.Publish(channel, execution)
	return nil
}
