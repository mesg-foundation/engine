package execution

import "fmt"

// StatusError is an error when the processing is done on en execution with the wrong status
type StatusError struct {
	ExpectedStatus Status
	ActualStatus   Status
}

// Error returns the string representation of error.
func (e StatusError) Error() string {
	return fmt.Sprintf("Execution status error: %q instead of %q", e.ActualStatus, e.ExpectedStatus)
}
