package api

import (
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// SubmitResult submits results for executionID.
// rerr used to submit result with an error.
func (a *API) SubmitResult(executionID string, outputKey string, outputData map[string]interface{}, rerr error) error {
	return newResultSubmitter(a).Submit(executionID, outputKey, outputData, rerr)
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
// rerr used to submit result with an error.
func (s *resultSubmitter) Submit(executionID string, outputKey string, outputData map[string]interface{}, rerr error) error {
	exec, err := s.api.execDB.Find(executionID)
	if err != nil {
		return err
	}
	exec.Service, err = service.FromService(exec.Service, service.ContainerOption(s.api.container))
	if err != nil {
		return err
	}

	if err = exec.Complete(outputKey, outputData, rerr); err != nil {
		exec.Error = err
	} else {
		if err = s.api.execDB.Save(exec); err != nil {
			exec.Error = err
		}
	}
	go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	return err
}
