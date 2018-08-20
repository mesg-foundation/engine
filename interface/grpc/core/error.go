package core

// NotRunningServiceError is an error when a service is not running.
type NotRunningServiceError struct {
	ServiceID string
}

func (e *NotRunningServiceError) Error() string {
	return "Service " + e.ServiceID + " is not running"
}
