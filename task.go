package mesg

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api/service"
	"golang.org/x/net/context"
)

// Task represents a MESG task.
type Task struct {
	name    string
	handler func(*Request)
}

// NewTask creates a task with name. handler executed when a matching task request received.
func NewTask(name string, handler func(*Request)) Task {
	t := Task{
		name:    name,
		handler: handler,
	}
	return t
}

// Request holds information about a Task request.
type Request struct {
	executionID string
	key         string
	data        string
	service     *Service
}

// Get task data input to out.
func (t *Request) Get(out interface{}) error {
	return json.Unmarshal([]byte(t.data), out)
}

// Reply sends data for key output as the response message of a Task.
func (t *Request) Reply(key string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = t.service.serviceClient.SubmitResult(context.Background(), &service.SubmitResultRequest{
		ExecutionID: t.executionID,
		OutputKey:   key,
		OutputData:  string(dataBytes),
	})
	return err
}
