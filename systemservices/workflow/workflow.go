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

// Delete stops and deletes workflow with id.
// TODO(ilgooz) close active log streams and delete the old log messages.
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
// TODO(ilgooz): support getting logs with workflow's name as well.
func (w *Workflow) Logs(id string) (stdLogs, errLogs io.ReadCloser, err error) {
	logs, err := w.api.ServiceLogs(w.serviceID, api.ServiceLogsDependenciesFilter("service"))
	if err != nil {
		return nil, nil, err
	}
	serviceLogs := logs[0]
	return newLogStream(id, serviceLogs.Standard), newLogStream(id, serviceLogs.Error), nil
}

// logStream filters service logs for getting logs of a workflow.
// it implements io.ReadCloser.
type logStream struct {
	workflowID string

	c io.Closer
	s *bufio.Scanner

	data []byte
	i    int64
}

// newLogStream returns a log stream that filters service logs for getting logs of workflowID.
func newLogStream(workflowID string, rc io.ReadCloser) *logStream {
	return &logStream{
		workflowID: workflowID,
		c:          rc,
		s:          bufio.NewScanner(rc),
	}
}

// logLine represents a log line received from log stream.
type logLine struct {
	WorkflowID string `json:"workflowID"`
}

// Read implements io.Reader.
func (s *logStream) Read(p []byte) (n int, err error) {
	if s.i >= int64(len(s.data)) {
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

func (s *logStream) scan() ([]byte, error) {
	if s.s.Scan() {
		data := s.s.Bytes()
		var line *logLine
		if err := json.Unmarshal(data, &line); err != nil {
			return nil, err
		}
		if line.WorkflowID == s.workflowID {
			return data, nil
		}
	}
	return nil, s.s.Err()
}

// Close implements io.Closer.
func (s *logStream) Close() error {
	return s.c.Close()
}
