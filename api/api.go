package api

import (
	"fmt"

	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	uuid "github.com/satori/go.uuid"
)

// API exposes all functionalities of MESG core.
type API struct {
	db     database.ServiceDB
	execDB database.ExecutionDB
	sm     service.Manager
}

// New creates a new API with given options.
func New(db database.ServiceDB, execDB database.ExecutionDB, sm service.Manager) *API {
	return &API{
		db:     db,
		execDB: execDB,
		sm:     sm,
	}
}

// GetService returns service serviceID.
func (api *API) GetService(service string) (*service.Service, error) {
	s, err := api.db.Get(service)
	if err != nil {
		return nil, err
	}

	if err := api.sm.Status(s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteService stops and deletes the service.
// If deleteData is set to true, any persistent data that belongs to
// service and its dependencies will be deleted.
func (api *API) DeleteService(service string, deleteData bool) error {
	s, err := api.db.Get(service)
	if err != nil {
		return err
	}

	if err := api.sm.Stop(s); err != nil {
		return err
	}
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := api.sm.Delete(s); err != nil {
			return err
		}
	}
	return api.db.Delete(service)
}

// ServiceLogs gives logs for all dependencies or one when specified with filters of service serviceID.
func (api *API) ServiceLogs(service string, dependencies []string) ([]*service.LogReader, error) {
	s, err := api.db.Get(service)
	if err != nil {
		return nil, err
	}
	return api.sm.Logs(s, dependencies)
}

// ListServices returns services matches with filters.
func (api *API) ListServices() ([]*service.Service, error) {
	services, err := api.db.All()
	if err != nil {
		return nil, err
	}

	for i := range services {
		if err := api.sm.Status(services[i]); err != nil {
			return nil, err
		}
	}
	return services, nil
}

// StopService stops service serviceID.
func (api *API) StopService(service string) error {
	s, err := api.db.Get(service)
	if err != nil {
		return err
	}
	return api.sm.Stop(s)
}

// StartService starts service serviceID.
func (api *API) StartService(service string) error {
	s, err := api.db.Get(service)
	if err != nil {
		return err
	}
	return api.sm.Start(s)
}

// EmitEvent emits a MESG event eventKey with eventData for service.
func (api *API) EmitEvent(service, eventKey string, eventData map[string]interface{}) error {
	s, err := api.db.Get(service)
	if err != nil {
		return err
	}

	if err := s.ValidateEventData(eventKey, eventData); err != nil {
		return err
	}

	e := event.New(s, eventKey, eventData)
	go pubsub.Publish(s.EventSubscriptionChannel(), e)
	return nil
}

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (api *API) ExecuteTask(serviceid, taskKey string, inputData map[string]interface{}, tags []string) (string, error) {
	s, err := api.db.Get(serviceid)
	if err != nil {
		return "", err
	}

	if err := api.sm.Status(s); err != nil {
		return "", err
	}

	if s.Status != service.StatusRunning {
		return "", fmt.Errorf("service %q is not running", s.Hash)
	}

	if err := s.ValidateTaskInputs(taskKey, inputData); err != nil {
		return "", err
	}

	eventID := uuid.NewV4().String()
	exec := execution.New(s, eventID, taskKey, inputData, tags)
	if err := exec.Execute(); err != nil {
		return "", err
	}

	if err := api.execDB.Save(exec); err != nil {
		return "", err
	}

	go pubsub.Publish(s.TaskSubscriptionChannel(), exec)
	return exec.ID, nil
}
