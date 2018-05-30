package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mesg-foundation/core/api/core"
)

func (task *Task) processEvent(data *core.EventData) (err error) {
	var d interface{}
	err = json.Unmarshal([]byte(data.EventData), &d)
	if err != nil {
		return
	}
	return task.process(d)
}

func (task *Task) processResult(data *core.ResultData) (err error) {
	var d interface{}
	err = json.Unmarshal([]byte(data.OutputData), &d)
	if err != nil {
		return
	}
	return task.process(d)
}

func (task *Task) process(data interface{}) (err error) {
	taskData := task.Inputs(data)
	var taskDataJSON []byte
	taskDataJSON, err = json.Marshal(taskData)
	if err != nil {
		return
	}
	log.Println("Trigger task", task.Name)
	_, err = getClient().ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: task.Service,
		TaskKey:   task.Name,
		TaskData:  string(taskDataJSON),
	})
	return
}
