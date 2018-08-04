package execution

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

// Execution stores all informations about executions.
type Execution struct {
	ID                string
	Service           *service.Service
	Task              string
	CreatedAt         time.Time
	ExecutedAt        time.Time
	ExecutionDuration time.Duration
	Inputs            map[string]interface{}
	Output            string
	OutputData        map[string]interface{}
}
