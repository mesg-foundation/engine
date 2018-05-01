package execution

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

type Execution struct {
	ID                string
	Service           *service.Service
	Task              string
	CreatedAt         time.Time
	ExecutedAt        time.Time
	ExecutionDuration time.Duration
	Inputs            interface{}
	Output            string
	OutputData        interface{}
}
