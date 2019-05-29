package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
)

// Execution holds information about a Task execution.
type Execution struct {
	// Hash is the execution id of task.
	Hash, hash string

	// Key is the name of task.
	Key string

	// inputs holds task inputs.
	inputs string

	service *Service
}

func newExecution(service *Service, data *serviceapi.TaskData) *Execution {
	return &Execution{
		Hash:    data.ExecutionHash,
		Key:     data.TaskKey,
		hash:    data.ExecutionHash,
		inputs:  data.InputData,
		service: service,
	}
}

// Data decodes task data input to out.
func (e *Execution) Data(out interface{}) error {
	return json.Unmarshal([]byte(e.inputs), out)
}

// reply sends task results to core.
func (e *Execution) reply(data interface{}, reterr error) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.service.callTimeout)
	defer cancel()
	if reterr != nil {
		_, err = e.service.client.SubmitResult(ctx, &serviceapi.SubmitResultRequest{
			ExecutionHash: e.hash,
			Result: &serviceapi.SubmitResultRequest_Error{
				Error: reterr.Error(),
			},
		})
	} else {
		resp, err1 := json.Marshal(data)
		if err1 != nil {
			return err1
		}
		_, err = e.service.client.SubmitResult(ctx, &serviceapi.SubmitResultRequest{
			ExecutionHash: e.hash,
			Result: &serviceapi.SubmitResultRequest_OutputData{
				OutputData: string(resp),
			},
		})
	}
	return err
}
