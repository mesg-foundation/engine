package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// SubmitResult submits results for executionID.
func (a *API) SubmitResult(executionID string, outputKey string, outputData interface{}) error {
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
func (s *resultSubmitter) Submit(executionID string, outputKey string, outputData interface{}) error {
	exec, err := s.api.execDB.Find(executionID)
	if err != nil {
		return err
	}
	exec.Service, err = service.FromService(exec.Service, service.ContainerOption(s.api.container))
	if err != nil {
		return err
	}

	// parse output data into a map.
	var outputDataMap map[string]interface{}
	switch d := outputData.(type) {
	case map[string]interface{}:
		outputDataMap = d
	case string:
		if err = json.Unmarshal([]byte(d), &outputDataMap); err != nil {
			err = fmt.Errorf("invalid output data error: %s", err)
		}
	default:
		err = errors.New("unexpected output data type")
	}

	if err == nil {
		err = exec.Complete(outputKey, outputDataMap)
	}

	if err != nil {
		exec.Failed(err)
	}

	if err = s.api.execDB.Save(exec); err != nil {
		return err
	}

	go pubsub.Publish(exec.Service.ResultSubscriptionChannel(), exec)
	return err
}
