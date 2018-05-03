package event

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

// Event is a type that store all informations about Events
type Event struct {
	Service   *service.Service
	Key       string
	Data      interface{}
	CreatedAt time.Time
}
