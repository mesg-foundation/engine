package api

import (
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
)

// taskExecutor provides functionalities to execute a MESG task.
type taskExecutor struct {
	api *API
}

// newTaskExecutor creates a new taskExecutor with given api.
func newTaskExecutor(api *API) *taskExecutor {
	return &taskExecutor{
		api: api,
	}
}

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (e *taskExecutor) Execute(serviceID, taskKey string, inputData map[string]interface{},
	tags []string) (executionID string, err error) {
	s, err := services.Get(serviceID)
	if err != nil {
		return "", err
	}
	if err := e.checkServiceStatus(&s); err != nil {
		return "", err
	}
	return e.execute(&s, taskKey, inputData, tags)
}

// checkServiceStatus checks service status. A task should be executed only if
// task's service is running.
func (e *taskExecutor) checkServiceStatus(s *service.Service) error {
	status, err := s.Status()
	if err != nil {
		return err
	}
	if status != service.RUNNING {
		return &NotRunningServiceError{ServiceID: s.Hash()}
	}
	return nil
}

// execute executes task.
func (e *taskExecutor) execute(s *service.Service, key string, inputs map[string]interface{},
	tags []string) (executionID string, err error) {
	exc, err := execution.Create(s, key, inputs, tags)
	if err != nil {
		return "", err
	}
	return exc.ID, exc.Execute()
}

// NotRunningServiceError is an error returned when the service is not running that
// a task needed to be executed on.
type NotRunningServiceError struct {
	ServiceID string
}

func (e *NotRunningServiceError) Error() string {
	return fmt.Sprintf("Service %q is not running", e.ServiceID)
}
