package api

import (
	"github.com/mesg-foundation/core/execution"
	uuid "github.com/satori/go.uuid"
)

// ExecuteAndListen executes given task and listenes for result.
func (a *API) ExecuteAndListen(serviceID, task string, inputs map[string]interface{}) (*execution.Execution, error) {
	tag := uuid.NewV4().String()
	_, err := a.ExecuteTask(serviceID, task, inputs, []string{tag})
	if err != nil {
		return nil, err

	}
	result, err := a.ListenResult(serviceID, ListenResultTagFilters([]string{tag}))
	if err != nil {
		return nil, err

	}
	defer result.Close()
	return <-result.Executions, nil

}
