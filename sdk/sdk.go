package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/execution"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
	"github.com/mesg-foundation/core/utils/hash"
	uuid "github.com/satori/go.uuid"
)

// SDK exposes all functionalities of MESG core.
type SDK struct {
	Service *servicesdk.Service

	ps *pubsub.PubSub

	m         manager.Manager
	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new SDK with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, execDB database.ExecutionDB) *SDK {
	return &SDK{
		Service:   servicesdk.New(m, c, db, execDB),
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
func (sdk *SDK) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{}, tags []string) (executionHash []byte, err error) {
	s, err := sdk.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	// a task should be executed only if task's service is running.
	status, err := sdk.m.Status(s)
	if err != nil {
		return nil, err
	}
	if status != service.RUNNING {
		return nil, &NotRunningServiceError{ServiceID: s.Sid}
	}

	if err := s.RequireTaskInputs(taskKey, inputData); err != nil {
		return nil, err
	}

	// execute the task.
	eventID := uuid.NewV4().String()
	exec := execution.New(s.Hash, nil, eventID, taskKey, inputData, tags)
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	if err = sdk.execDB.Save(exec); err != nil {
		return nil, err
	}
	go sdk.ps.Pub(exec, executionSubTopic(s.Hash))
	return exec.Hash, nil
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
func (sdk *SDK) SubmitResult(executionHash []byte, outputs []byte, reterr error) error {
	exec, err := sdk.processExecution(executionHash, outputs, reterr)
	if err != nil {
		return err
	}

	go sdk.ps.Pub(exec, executionSubTopic(exec.ServiceHash))
	return nil
}

// processExecution processes execution and marks it as complated or failed.
func (sdk *SDK) processExecution(executionHash []byte, outputData []byte, reterr error) (*execution.Execution, error) {
	tx, err := sdk.execDB.OpenTransaction()
	if err != nil {
		return nil, err
	}

	exec, err := tx.Find(executionHash)
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
		o, err := sdk.validateExecutionOutput(exec.ServiceHash, exec.TaskKey, outputData)
		if err != nil {
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

func (sdk *SDK) validateExecutionOutput(service, taskKey string, jsonout []byte) (map[string]interface{}, error) {
	var output map[string]interface{}
	if err := json.Unmarshal(jsonout, &output); err != nil {
		return nil, fmt.Errorf("invalid output: %s", err)
	}

	s, err := sdk.db.Get(service)
	if err != nil {
		return nil, err
	}

	if err := s.RequireTaskOutputs(taskKey, output); err != nil {
		return nil, err
	}
	return output, nil
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
