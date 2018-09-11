package container

// StatusType of the service.
type StatusType uint

// Possible status for services.
const (
	UNKNOWN StatusType = iota
	STOPPED
	STARTING
	RUNNING
)
