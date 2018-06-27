package client

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

func (task *Task) processEvent(wf *Workflow, data *core.EventData) (err error) {
	var d interface{}
	err = json.Unmarshal([]byte(data.EventData), &d)
	if err != nil {
		return
	}
	return task.process(wf, d)
}

func (task *Task) processResult(wf *Workflow, data *core.ResultData) (err error) {
	var d interface{}
	err = json.Unmarshal([]byte(data.OutputData), &d)
	if err != nil {
		return
	}
	return task.process(wf, d)
}

func (task *Task) process(wf *Workflow, data interface{}) (err error) {
	inputData, err := task.convertData(data)
	if err != nil {
		return
	}
	_, err = wf.client.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: task.ServiceID,
		TaskKey:   task.Name,
		InputData: inputData,
	})
	return
}

func (task *Task) convertData(data interface{}) (res string, err error) {
	inputData := task.Inputs(data)
	var inputDataJSON []byte
	inputDataJSON, err = json.Marshal(inputData)
	if err != nil {
		return
	}
	res = string(inputDataJSON)
	return
}
