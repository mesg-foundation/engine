package container

import (
	"time"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 1
	RUNNING StatusType = 2
)

// TimeoutError represents an error of timeout
type TimeoutError struct {
	duration time.Duration
	name     string
}

func (e *TimeoutError) Error() string {
	return "Timeout reached after " + e.duration.String() + " for ressource " + e.name
}
