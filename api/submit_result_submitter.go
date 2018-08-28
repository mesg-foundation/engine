package api

import (
	"fmt"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
)

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

	serviceOutput, outputFound := execution.Service.Tasks[execution.Task].Outputs[outputKey]
	if !outputFound {
		return &service.TaskOutputNotFoundError{
			TaskKey:       execution.Task,
			TaskOutputKey: outputKey,
			ServiceName:   execution.Service.Name,
		}
	}
	warnings := execution.Service.ValidateParametersSchema(serviceOutput.Data, outputData)
	if len(warnings) > 0 {
		return &service.InvalidTaskOutputError{
			TaskKey:       execution.Task,
			TaskOutputKey: outputKey,
			Warnings:      warnings,
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
