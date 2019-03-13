package api

import (
	"time"

	"github.com/mesg-foundation/core/execution"
	uuid "github.com/satori/go.uuid"
)

// ExecuteAndListen executes given task and listen for result.
func (a *API) ExecuteAndListen(serviceID, task string, inputs map[string]interface{}) (*execution.Execution, error) {
	tag := uuid.NewV4().String()
	result, err := a.ListenResult(serviceID, ListenResultTagFilters([]string{tag}))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	// XXX: sleep because listen stream may not be ready to stream the data
	// and execution will done before stream is ready. In that case the response
	// wlll never come TODO: investigate
	time.Sleep(1 * time.Second)

	if _, err := a.ExecuteTask(serviceID, task, inputs, []string{tag}); err != nil {
		return nil, err
	}
	return <-result.Executions, nil

}
