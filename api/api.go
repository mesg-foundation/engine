package api

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

// API exposes all functionalities of MESG core.
type API struct {
	ps *pubsub.PubSub

	m         manager.Manager
	container container.Container
	db        database.ServiceDB
	execDB    database.ExecutionDB
}

// New creates a new API with given options.
func New(m manager.Manager, c container.Container, db database.ServiceDB, execDB database.ExecutionDB) *API {
	return &API{
		ps:        pubsub.New(0),
		m:         m,
		container: c,
		db:        db,
		execDB:    execDB,
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
	return a.m.Status(service)
}

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	_, err = a.m.Start(s)
	return err
}

// StopService stops service serviceID.
func (a *API) StopService(serviceID string) error {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	return a.m.Stop(s)
}

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// the service and to its dependencies also will be deleted.
func (a *API) DeleteService(serviceID string, deleteData bool) error {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	if err := a.m.Stop(s); err != nil {
		return err
	}
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := a.m.Delete(s); err != nil {
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
	e, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}

	go a.ps.Pub(e, eventSubTopic(s.Hash))
	return nil
}

// ExecuteTask executes a task tasKey with inputData and tags for service serviceID.
func (a *API) ExecuteTask(serviceID, taskKey string, inputData map[string]interface{}, tags []string) (executionHash []byte, err error) {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	// a task should be executed only if task's service is running.
	status, err := a.m.Status(s)
	if err != nil {
		return nil, err
	}
	if status != service.RUNNING {
		return nil, &NotRunningServiceError{ServiceID: s.Sid}
	}

	task, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	if err := task.RequireInputs(inputData); err != nil {
		return nil, err
	}

	// execute the task.
	eventID := uuid.NewV4().String()
	exec := execution.New(s.Hash, nil, eventID, taskKey, inputData, tags)
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	if err = a.execDB.Save(exec); err != nil {
		return nil, err
	}
	go a.ps.Pub(exec, executionSubTopic(s.Hash))
	return exec.Hash, nil
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (a *API) ListenEvent(service string, f *EventFilter) (*EventListener, error) {
	s, err := a.db.Get(service)
	if err != nil {
		return nil, err
	}

	if f.HasKey() {
		if _, err := s.GetEvent(f.Key); err != nil {
			return nil, err
		}
	}

	l := NewEventListener(a.ps, eventSubTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// ListenExecution listens executions on service.
func (a *API) ListenExecution(service string, f *ExecutionFilter) (*ExecutionListener, error) {
	s, err := a.db.Get(service)
	if err != nil {
		return nil, err
	}

	if f != nil {
		if f.TaskKey == "" && f.OutputKey != "" {
			return nil, fmt.Errorf("execution filter: output key given without task key")
		}
		if f.HasTaskKey() {
			task, err := s.GetTask(f.TaskKey)
			if err != nil {
				return nil, err
			}
			if f.HasOutputKey() {
				if _, err := task.GetOutput(f.OutputKey); err != nil {
					return nil, err
				}
			}
		}
	}

	l := NewExecutionListener(a.ps, executionSubTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// SubmitResult submits results for executionHash.
func (a *API) SubmitResult(executionHash []byte, outputKey string, outputs []byte) error {
	exec, stateChanged, err := a.processExecution(executionHash, outputKey, outputs)
	if stateChanged {
		// only publish to listeners when the execution's state changed.
		go a.ps.Pub(exec, executionSubTopic(exec.ServiceHash))
	}
	return err
}

// processExecution processes execution and marks it as complated or failed.
func (a *API) processExecution(executionHash []byte, outputKey string, outputData []byte) (exec *execution.Execution, stateChanged bool, err error) {
	stateChanged = false
	tx, err := a.execDB.OpenTransaction()
	if err != nil {
		return nil, false, err
	}

	exec, err = tx.Find(executionHash)
	if err != nil {
		tx.Discard()
		return nil, false, err
	}

	output, err := a.validateExecutionOutput(exec.ServiceHash, exec.TaskKey, outputKey, outputData)
	if err != nil {
		return a.saveExecution(tx, exec, err)
	}

	if err := exec.Complete(outputKey, output); err != nil {
		return a.saveExecution(tx, exec, err)
	}

	return a.saveExecution(tx, exec, nil)
}

func (a *API) saveExecution(tx database.ExecutionTransaction, exec *execution.Execution, err error) (execOut *execution.Execution, stateChanged bool, errOut error) {
	if err != nil {
		if errFailed := exec.Failed(err); errFailed != nil {
			tx.Discard()
			return exec, false, errFailed
		}
	}
	if errSave := tx.Save(exec); errSave != nil {
		tx.Discard()
		return exec, true, errSave
	}
	if errCommit := tx.Commit(); errCommit != nil {
		return exec, true, errCommit
	}
	return exec, true, err
}

func (a *API) validateExecutionOutput(service, taskKey, outputKey string, jsonout []byte) (map[string]interface{}, error) {
	var output map[string]interface{}
	if err := json.Unmarshal(jsonout, &output); err != nil {
		return nil, fmt.Errorf("invalid output: %s", err)
	}

	s, err := a.db.Get(service)
	if err != nil {
		return nil, err
	}

	task, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	o, err := task.GetOutput(outputKey)
	if err != nil {
		return nil, err
	}
	return output, o.RequireData(output)
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
