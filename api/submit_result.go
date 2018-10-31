package api

import (
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
	exec, err := s.api.execDB.Find(executionID)
	if err != nil {
		return err
	}
	if err := exec.Complete(outputKey, outputData); err != nil {
		return err
	}
	if err = s.api.execDB.Save(exec); err != nil {
		return err
	}
	go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	return nil
}
