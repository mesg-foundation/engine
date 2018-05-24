package docker

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 1
	RUNNING StatusType = 2
)
