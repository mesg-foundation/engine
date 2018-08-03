package container

// StatusType of the service.
type StatusType uint

// Possible status for services.
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
)
