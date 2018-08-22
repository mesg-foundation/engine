package mesg

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/service"
)

// Execution holds information about a Task execution.
type Execution struct {
	// ID is the execution id of task.
	ID, id string

	// Key is the name of task.
	Key string

	// inputs holds task inputs.
	inputs string

	service *Service
}

func newExecution(service *Service, data *service.TaskData) *Execution {
	return &Execution{
		ID:      data.ExecutionID,
		Key:     data.TaskKey,
		id:      data.ExecutionID,
		inputs:  data.InputData,
		service: service,
	}
}

// Data decodes task data input to out.
func (e *Execution) Data(out interface{}) error {
	return json.Unmarshal([]byte(e.inputs), out)
}

// reply sends task results to core.
func (e *Execution) reply(key string, data Data) error {
	if err := e.validateTaskOutputs(key, data); err != nil {
		return err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.service.callTimeout)
	defer cancel()
	_, err = e.service.client.SubmitResult(ctx, &service.SubmitResultRequest{
		ExecutionID: e.id,
		OutputKey:   string(key),
		OutputData:  string(dataBytes),
	})
	return err
}

// validateTaskOutputs validates output key and data of task as described in mesg.yaml.
// TODO(ilgooz) use validation handlers of core server to do this?
func (e *Execution) validateTaskOutputs(key string, data Data) error { return nil }

type errTaskOutput struct{}

func (e errTaskOutput) Error() string {
	return ""
}
