package execution

// StatusError is an error when the processing is done on en execution with the wrong status
type StatusError struct{}

// Error returns the string representation of error.
func (e StatusError) Error() string {
	return "Execution status error"
}
