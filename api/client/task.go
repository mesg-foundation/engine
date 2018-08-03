package client

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

func (task *Task) processEvent(wf *Workflow, data *core.EventData) error {
	var d interface{}
	if err := json.Unmarshal([]byte(data.EventData), &d); err != nil {
		return err
	}
	return task.process(wf, d)
}

func (task *Task) processResult(wf *Workflow, data *core.ResultData) error {
	var d interface{}
	if err := json.Unmarshal([]byte(data.OutputData), &d); err != nil {
		return err
	}
	return task.process(wf, d)
}

func (task *Task) process(wf *Workflow, data interface{}) error {
	inputData, err := task.convertData(data)
	if err != nil {
		return err
	}
	_, err = wf.client.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: task.ServiceID,
		TaskKey:   task.Name,
		InputData: inputData,
	})
	return err
}

func (task *Task) convertData(data interface{}) (string, error) {
	inputData := task.Inputs(data)
	inputDataJSON, err := json.Marshal(inputData)
	if err != nil {
		return "", err
	}
	return string(inputDataJSON), nil
}
