package api

import (
	"fmt"

	"github.com/mesg-foundation/core/execution"
)

// SubmitResult submits results for executionID.
func (a *API) SubmitResult(executionID, outputKey string, outputData map[string]interface{}) error {
	return newResultSubmitter(a).Submit(executionID, outputKey, outputData)
}

// resultSubmitter provides functionalities to submit a MESG task result.
type resultSubmitter struct {
	api *API
}

// newResultSubmitter creates a new resultSubmitter with given api.
func newResultSubmitter(api *API) *resultSubmitter {
	return &resultSubmitter{
		api: api,
	}
}

// Submit submits results for executionID.
func (s *resultSubmitter) Submit(executionID, outputKey string, outputData map[string]interface{}) error {
	execution := execution.InProgress(executionID)
	if execution == nil {
		return &MissingExecutionError{
			ID: executionID,
		}
	}
	return execution.Complete(outputKey, outputData)
}

// MissingExecutionError is returned when corresponding execution doesn't exists.
type MissingExecutionError struct {
	ID string
}

func (e *MissingExecutionError) Error() string {
	return fmt.Sprintf("Execution %q doesn't exists", e.ID)
}
