package execution

import (
	"log"
	"time"
)

func (execution *Execution) Complete(output string, data interface{}) (err error) {
	err = execution.moveFromInProgressToProcessed()
	if err != nil {
		return
	}
	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = output
	execution.OutputData = data
	log.Println("[COMPLETED]", execution.Task)
	return
}
