package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
	"github.com/mesg-foundation/core/utils/hash"
	uuid "github.com/satori/go.uuid"
)

// SDK exposes all functionalities of MESG core.
type SDK struct {
	ps *pubsub.PubSub

	m         manager.Manager
	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new SDK with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, execDB database.ExecutionDB) *SDK {
	return &SDK{
		ps:        pubsub.New(0),
		m:         m,
		container: c,
		db:        db,
		execDB:    execDB,
	}
}

// GetService returns service serviceID.
func (sdk *SDK) GetService(serviceID string) (*service.Service, error) {
	return sdk.db.Get(serviceID)
}

// ListServices returns all services.
func (sdk *SDK) ListServices() ([]*service.Service, error) {
	return sdk.db.All()
}

// Status returns the status of a service
func (sdk *SDK) Status(service *service.Service) (service.StatusType, error) {
	return sdk.m.Status(service)
}

// StartService starts service serviceID.
func (sdk *SDK) StartService(serviceID string) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}
	_, err = sdk.m.Start(s)
	return err
}

// StopService stops service serviceID.
func (sdk *SDK) StopService(serviceID string) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}
	return sdk.m.Stop(s)
}

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// the service and to its dependencies also will be deleted.
func (sdk *SDK) DeleteService(serviceID string, deleteData bool) error {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return err
	}
	if err := sdk.m.Stop(s); err != nil {
		return err
	}
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := sdk.m.Delete(s); err != nil {
			return err
		}
	}
	return sdk.db.Delete(serviceID)
}

// EmitEvent emits a MESG event eventKey with eventData for service token.
func (sdk *SDK) EmitEvent(token, eventKey string, eventData map[string]interface{}) error {
	s, err := sdk.db.Get(token)
	if err != nil {
		return err
	}
	e, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}

	go sdk.ps.Pub(e, eventSubTopic(s.Hash))
	return nil
}

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (sdk *SDK) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{},
	tags []string) (executionID string, err error) {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return "", err
	}
	// a task should be executed only if task's service is running.
	status, err := sdk.m.Status(s)
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
	if err = sdk.execDB.Save(exec); err != nil {
		return "", err
	}
	go sdk.ps.Pub(exec, executionSubTopic(s.Hash))
	return exec.ID, nil
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (sdk *SDK) ListenEvent(service string, f *EventFilter) (*EventListener, error) {
	s, err := sdk.db.Get(service)
	if err != nil {
		return nil, err
	}

	if f.HasKey() {
		if _, err := s.GetEvent(f.Key); err != nil {
			return nil, err
		}
	}

	l := NewEventListener(sdk.ps, eventSubTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// ListenExecution listens executions on service.
func (sdk *SDK) ListenExecution(service string, f *ExecutionFilter) (*ExecutionListener, error) {
	s, err := sdk.db.Get(service)
	if err != nil {
		return nil, err
	}

	if f != nil && f.HasTaskKey() {
		if _, err := s.GetTask(f.TaskKey); err != nil {
			return nil, err
		}
	}

	l := NewExecutionListener(sdk.ps, executionSubTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// SubmitResult submits results for executionID.
func (sdk *SDK) SubmitResult(executionID string, outputs []byte, reterr error) error {
	exec, err := sdk.processExecution(executionID, outputs, reterr)
	if err != nil {
		return err
	}

	go sdk.ps.Pub(exec, executionSubTopic(exec.Service.Hash))
	return nil
}

// processExecution processes execution and marks it as complated or failed.
func (sdk *SDK) processExecution(executionID string, outputData []byte, reterr error) (*execution.Execution, error) {
	tx, err := sdk.execDB.OpenTransaction()
	if err != nil {
		return nil, err
	}

	exec, err := tx.Find(executionID)
	if err != nil {
		tx.Discard()
		return nil, err
	}

	if reterr != nil {
		if err := exec.Failed(reterr); err != nil {
			tx.Discard()
			return nil, err
		}
	} else {
		var o map[string]interface{}
		if err := json.Unmarshal(outputData, &o); err != nil {
			return nil, err
		}

		if err := exec.Complete(o); err != nil {
			return nil, err
		}
	}

	if err := tx.Save(exec); err != nil {
		tx.Discard()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Discard()
		return nil, err
	}

	return exec, nil
}

// NotRunningServiceError is an error returned when the service is not running that
// a task needed to be executed on.
type NotRunningServiceError struct {
	ServiceID string
}

func (e *NotRunningServiceError) Error() string {
	return fmt.Sprintf("Service %q is not running", e.ServiceID)
}

const (
	eventTopic     = "Event"
	executionTopic = "Execution"
)

// eventSubTopic returns the topic to listen for events from this service.
func eventSubTopic(serviceHash string) string {
	return hash.Calculate([]string{serviceHash, eventTopic})
}

// executionSubTopic returns the topic to listen for tasks from this service.
func executionSubTopic(serviceHash string) string {
	return hash.Calculate([]string{serviceHash, executionTopic})
}
