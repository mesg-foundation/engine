package executionsdk

import (
	"encoding/json"
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/mr-tron/base58"
	uuid "github.com/satori/go.uuid"
)

const (
	// streamTopic is topic used to broadcast executions.
	streamTopic = "execution-stream"
	topic       = "Execution"
)

// Execution exposes execution APIs of MESG.
type Execution struct {
	ps     *pubsub.PubSub
	m      manager.Manager
	db     database.ServiceDB
	execDB database.ExecutionDB
}

// New creates a new Execution SDK with given options.
func New(m manager.Manager, ps *pubsub.PubSub, db database.ServiceDB, execDB database.ExecutionDB) *Execution {
	return &Execution{
		m:      m,
		ps:     ps,
		db:     db,
		execDB: execDB,
	}
}

// Get returns execution that matches given hash.
func (e *Execution) Get(hash []byte) (*execution.Execution, error) {
	return e.execDB.Find(hash)
}

// GetStream returns execution that matches given hash.
func (e *Execution) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}

// Update updates execution that matches given hash.
func (e *Execution) Update(executionHash, outputs []byte, reterr error) error {
	exec, err := e.processExecution(executionHash, outputs, reterr)
	if err != nil {
		return err
	}

	go e.ps.Pub(exec, streamTopic)
	go e.ps.Pub(exec, subTopic(exec.ServiceHash))
	return nil
}

// processExecution processes execution and marks it as complated or failed.
func (e *Execution) processExecution(executionHash []byte, outputData []byte, reterr error) (*execution.Execution, error) {
	tx, err := e.execDB.OpenTransaction()
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
		o, err := e.validateExecutionOutput(exec.ServiceHash, exec.TaskKey, outputData)
		if err != nil {
			tx.Discard()
			return nil, err
		}

		if err := exec.Complete(o); err != nil {
			tx.Discard()
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

func (e *Execution) validateExecutionOutput(serviceHash []byte, taskKey string, jsonout []byte) (map[string]interface{}, error) {
	var output map[string]interface{}
	if err := json.Unmarshal(jsonout, &output); err != nil {
		return nil, fmt.Errorf("invalid output: %s", err)
	}

	s, err := e.db.Get(serviceHash)
	if err != nil {
		return nil, err
	}

	if err := s.RequireTaskOutputs(taskKey, output); err != nil {
		return nil, err
	}
	return output, nil
}

// Execute executes a task tasKey with inputData and tags for service serviceID.
func (e *Execution) Execute(serviceHash []byte, taskKey string, inputData map[string]interface{}, tags []string) (executionHash []byte, err error) {
	s, err := e.db.Get(serviceHash)
	if err != nil {
		return nil, err
	}
	// a task should be executed only if task's service is running.
	status, err := e.m.Status(s)
	if err != nil {
		return nil, err
	}
	if status != service.RUNNING {
		return nil, &NotRunningServiceError{ServiceID: base58.Encode(s.Hash)}
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
	if err = e.execDB.Save(exec); err != nil {
		return nil, err
	}

	go e.ps.Pub(exec, streamTopic)
	go e.ps.Pub(exec, subTopic(s.Hash))
	return exec.Hash, nil
}

// Listen listens executions on service.
func (e *Execution) Listen(serviceHash []byte, f *Filter) (*Listener, error) {
	s, err := e.db.Get(serviceHash)
	if err != nil {
		return nil, err
	}

	if f != nil && f.HasTaskKey() {
		if _, err := s.GetTask(f.TaskKey); err != nil {
			return nil, err
		}
	}

	l := NewListener(e.ps, subTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// subTopic returns the topic to listen for tasks from this service.
func subTopic(serviceHash []byte) string {
	return hash.Calculate([]string{base58.Encode(serviceHash), topic})
}

// NotRunningServiceError is an error returned when the service is not running that
// a task needed to be executed on.
type NotRunningServiceError struct {
	ServiceID string
}

func (e *NotRunningServiceError) Error() string {
	return fmt.Sprintf("Service %q is not running", e.ServiceID)
}
