package api

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// SubmitResult submits results for executionID.
func (a *API) SubmitResult(executionID string, outputKey string, outputData []byte) error {
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
func (s *resultSubmitter) Submit(executionID string, outputKey string, outputData []byte) error {
	exec, err := s.processExecution(executionID, outputKey, outputData)
	go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	return err
}

// processExecution processes execution and marks it as complated or failed.
func (s *resultSubmitter) processExecution(executionID string, outputKey string, outputData []byte) (*execution.Execution, error) {
	exec, err := s.api.execDB.Find(executionID)
	if err != nil {
		return nil, err
	}

	exec.Service, err = service.FromService(exec.Service, service.ContainerOption(s.api.container))
	if err != nil {
		exec.Failed(err)
		return exec, err
	}

	var outputDataMap map[string]interface{}
	if err = json.Unmarshal(outputData, &outputDataMap); err != nil {
		err = fmt.Errorf("invalid output data error: %s", err)
		exec.Failed(err)
	} else {
		if err = exec.Complete(outputKey, outputDataMap); err != nil {
			exec.Failed(err)
		}
	}

	if err := s.api.execDB.Save(exec); err != nil {
		exec.Failed(err)
		return exec, err
	}

	return exec, nil
}
