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
	exec, stateChanged, err := s.processExecution(executionID, outputKey, outputData)
	if stateChanged {
		// only publish to listeners when the execution's state changed.
		go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	}
	// always return any error to the service.
	return err
}

// processExecution processes execution and marks it as complated or failed.
func (s *resultSubmitter) processExecution(executionID string, outputKey string, outputData []byte) (exec *execution.Execution, stateChanged bool, err error) {
	stateChanged = false
	exec, err = s.api.execDB.Find(executionID)
	if err != nil {
		return nil, false, err
	}

	exec.Service, err = service.FromService(exec.Service, service.ContainerOption(s.api.container))
	if err != nil {
		return s.saveExecution(exec, err)
	}

	var outputDataMap map[string]interface{}
	if err := json.Unmarshal(outputData, &outputDataMap); err != nil {
		return s.saveExecution(exec, fmt.Errorf("invalid output data error: %s", err))
	}

	if err := exec.Complete(outputKey, outputDataMap); err != nil {
		return s.saveExecution(exec, err)
	}

	return s.saveExecution(exec, nil)
}

func (s *resultSubmitter) saveExecution(exec *execution.Execution, err error) (execOut *execution.Execution, stateChanged bool, errOut error) {
	if err != nil {
		if errFailed := exec.Failed(err); errFailed != nil {
			return exec, false, errFailed
		}
	}
	if errSave := s.api.execDB.Save(exec); errSave != nil {
		return exec, true, errSave
	}
	return exec, true, err
}
