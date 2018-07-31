package mesg

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/api/service"
)

var errOnlyOneOutputKey = errors.New("there should be only one output key in response")

// Task represents a MESG task.
type Task struct {
	name    string
	handler func(*Request) Response
}

// NewTask creates a task with name, handler executed when a matching task request is received.
func NewTask(name string, handler func(*Request) Response) Task {
	t := Task{
		name:    name,
		handler: handler,
	}
	return t
}

// Key is the output key of task.
type Key string

// Data is an event data or output data of a task.
type Data interface{}

// Response is a task response.
type Response map[Key]Data

// Request holds information about a Task request.
type Request struct {
	// ExecutionID is the execution id of task.
	ExecutionID string

	// Key is the name of task.
	Key string

	data    string
	service *Service
}

// Decode decodes task data input to out.
func (t *Request) Decode(out interface{}) error {
	return json.Unmarshal([]byte(t.data), out)
}

func (t *Request) reply(resp Response) error {
	if len(resp) != 1 {
		return errOnlyOneOutputKey
	}
	for key, data := range resp {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = t.service.client.SubmitResult(context.Background(), &service.SubmitResultRequest{
			ExecutionID: t.ExecutionID,
			OutputKey:   string(key),
			OutputData:  string(dataBytes),
		})
		return err
	}
	return nil
}
