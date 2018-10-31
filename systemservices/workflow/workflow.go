package workflow

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/mesg-foundation/core/api"
)

// WSS's tasks.
const (
	createTaskKey = "create"
	getTaskKey    = "get"
	deleteTaskKey = "delete"
)

// Workflow is a high level wrapper for Workflow System Service.
// It calls the WSS's tasks and reacts to its event through network.
// WSS responsible for managing and running workflows.
type Workflow struct {
	api       *api.API
	serviceID string
}

// New creates a new Workflow for given WSS serviceID and api.
func New(serviceID string, api *api.API) *Workflow {
	return &Workflow{
		api:       api,
		serviceID: serviceID,
	}
}

// Create creates and runs a workflow file with an optionally given unique name.
func (w *Workflow) Create(file []byte, name string) (id string, err error) {
	e, err := w.api.ExecuteAndListen(w.serviceID, createTaskKey, map[string]interface{}{
		"file": string(file),
		"name": name,
	})
	if err != nil {
		return "", err
	}

	switch e.Output {
	case "success":
		return e.OutputData["id"].(string), nil
	case "error":
		return "", errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// WorkflowDocument keeps workflow info.
type WorkflowDocument struct {
	// ID is the unique id for workflow.
	ID string

	// CreationID is the unique random id generated when
	// workflow is created.
	CreationID string

	// Name is the optionally set unique name for workflow.
	Name string
}

// Get returns the workflow info.
func (w *Workflow) Get(id string) (*WorkflowDocument, error) {
	e, err := w.api.ExecuteAndListen(w.serviceID, getTaskKey, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	switch e.Output {
	case "success":
		workflow := e.OutputData["workflow"].(map[string]interface{})
		return &WorkflowDocument{
			ID:         workflow["id"].(string),
			CreationID: workflow["creationID"].(string),
			Name:       workflow["name"].(string),
		}, nil
	case "error":
		return nil, errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// Delete stops and deletes workflow with id.
func (w *Workflow) Delete(id string) (err error) {
	e, err := w.api.ExecuteAndListen(w.serviceID, deleteTaskKey, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	switch e.Output {
	case "success":
		return nil
	case "error":
		return errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// Logs returns the standard and error log streams of workflow with id.
func (w *Workflow) Logs(id string) (stdLogs, errLogs io.ReadCloser, err error) {
	wdoc, err := w.Get(id)
	if err != nil {
		return nil, nil, err
	}
	logs, err := w.api.ServiceLogs(w.serviceID, api.ServiceLogsDependenciesFilter("service"))
	if err != nil {
		return nil, nil, err
	}
	serviceLogs := logs[0]
	stdout := newLogStream(wdoc.CreationID, serviceLogs.Standard)
	stderr := newLogStream(wdoc.CreationID, serviceLogs.Error)
	return stdout, stderr, nil
}

// logStream filters service logs for getting logs of a workflow.
// it implements io.ReadCloser.
type logStream struct {
	workflowCreationID string

	closing bool

	c io.Closer
	s *bufio.Scanner

	data []byte
	i    int64
}

// newLogStream returns a log stream that filters service logs to get logs
// of a workflow with workflowCreationID.
func newLogStream(workflowCreationID string, rc io.ReadCloser) *logStream {
	return &logStream{
		workflowCreationID: workflowCreationID,
		c:                  rc,
		s:                  bufio.NewScanner(rc),
	}
}

// logLine represents a log line received from log stream.
type logLine struct {
	WorkflowCreationID string `json:"workflowCreationID"`
	Workflow           struct {
		Deleted bool `json:"deleted"`
	} `json:"workflow"`
}

// Read implements io.Reader.
func (s *logStream) Read(p []byte) (n int, err error) {
	if s.i >= int64(len(s.data)) {
		if s.closing {
			return 0, io.EOF
		}
		data, err := s.scan()
		if err != nil {
			return 0, err
		}
		s.data = data
		s.i = 0
		return s.Read(p)
	}
	n = copy(p, s.data[s.i:])
	s.i += int64(n)
	return n, nil
}

// scan scans the next log line and compares the workflow ids to filter.
// it returns the line if ids matches.
func (s *logStream) scan() ([]byte, error) {
	// TODO(ilgooz)
	// read by line is error prone because json message can have unescaped \n char.
	if s.s.Scan() {
		data := s.s.Bytes()
		var line *logLine
		if err := json.Unmarshal(data, &line); err != nil {
			return nil, err
		}
		if line.WorkflowCreationID == s.workflowCreationID {
			if line.Workflow.Deleted {
				s.closing = true
			}
			return data, nil
		}
	}
	return nil, s.s.Err()
}

// Close implements io.Closer.
func (s *logStream) Close() error {
	return s.c.Close()
}
