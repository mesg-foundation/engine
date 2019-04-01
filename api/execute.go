package api

import (
	"fmt"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	uuid "github.com/satori/go.uuid"
)

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (a *API) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{},
	tags []string) (executionID string, err error) {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return "", err
	}
	s, err = service.FromService(s, service.ContainerOption(a.container))
	if err != nil {
		return "", err
	}

	// a task should be executed only if task's service is running.
	status, err := s.Status()
	if err != nil {
		return "", err
	}
	if status != service.RUNNING {
		return "", &NotRunningServiceError{ServiceID: s.Sid}
	}

	// execute the task.
	eventID := uuid.NewV4().String()
	exec, err := execution.New(s, eventID, taskKey, inputData, tags)
	if err != nil {
		return "", err
	}
	if err := exec.Execute(); err != nil {
		return "", err
	}
	if err = a.execDB.Save(exec); err != nil {
		return "", err
	}
	go pubsub.Publish(s.TaskSubscriptionChannel(), exec)
	return exec.ID, nil
}

// NotRunningServiceError is an error returned when the service is not running that
// a task needed to be executed on.
type NotRunningServiceError struct {
	ServiceID string
}

func (e *NotRunningServiceError) Error() string {
	return fmt.Sprintf("Service %q is not running", e.ServiceID)
}
