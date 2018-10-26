package api

import (
	"fmt"

	"github.com/mesg-foundation/core/pubsub"
)

// SubmitResult submits results for executionID.
func (a *API) SubmitResult(executionID string, outputKey string, outputData map[string]interface{}) error {
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
func (s *resultSubmitter) Submit(executionID string, outputKey string, outputData map[string]interface{}) error {
	execution, err := s.api.execDB.Complete(executionID, outputKey, outputData)
	if err != nil {
		return err
	}
	srv, err := s.api.db.Get(execution.ServiceID)
	if err != nil {
		return err
	}
	go pubsub.Publish(srv.ResultSubscriptionChannel(), execution)
	return nil
}

// MissingExecutionError is returned when corresponding execution doesn't exists.
type MissingExecutionError struct {
	ID string
}

func (e *MissingExecutionError) Error() string {
	return fmt.Sprintf("Execution %q doesn't exists", e.ID)
}
