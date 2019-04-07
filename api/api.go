package api

import (
	"fmt"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	uuid "github.com/satori/go.uuid"
)

// API exposes all functionalities of MESG core.
type API struct {
	db        database.ServiceDB
	execDB    database.ExecutionDB
	container container.Container
}

// Option is a configuration func for MESG.
type Option func(*API)

// New creates a new API with given options.
func New(db database.ServiceDB, execDB database.ExecutionDB, options ...Option) (*API, error) {
	a := &API{db: db, execDB: execDB}
	for _, option := range options {
		option(a)
	}
	if a.container == nil {
		var err error
		a.container, err = container.New()
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// ContainerOption configures underlying container access API.
func ContainerOption(container container.Container) Option {
	return func(a *API) {
		a.container = container
	}
}

// GetService returns service serviceID.
func (a *API) GetService(serviceID string) (*service.Service, error) {
	return a.db.Get(serviceID)
}

// ListServices returns all services.
func (a *API) ListServices() ([]*service.Service, error) {
	return a.db.All()
}

// Status returns the status of a service
func (a *API) Status(service *service.Service) (service.StatusType, error) {
	return service.Status(a.container)
}

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	sr, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	_, err = sr.Start(a.container)
	return err
}

// StopService stops service serviceID.
func (a *API) StopService(serviceID string) error {
	sr, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	return sr.Stop(a.container)
}

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// the service and to its dependencies also will be deleted.
func (a *API) DeleteService(serviceID string, deleteData bool) error {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	if err := s.Stop(a.container); err != nil {
		return err
	}
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := s.DeleteVolumes(a.container); err != nil {
			return err
		}
	}
	return a.db.Delete(serviceID)
}

// EmitEvent emits a MESG event eventKey with eventData for service token.
func (a *API) EmitEvent(token, eventKey string, eventData map[string]interface{}) error {
	s, err := a.db.Get(token)
	if err != nil {
		return err
	}
	event, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}
	event.Publish()
	return nil
}

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (a *API) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{},
	tags []string) (executionID string, err error) {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return "", err
	}
	// a task should be executed only if task's service is running.
	status, err := s.Status(a.container)
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
