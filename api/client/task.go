package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mesg-foundation/core/api/core"
)

func (task *Task) processEvent(client core.CoreClient, data *core.EventData) (err error) {
	var d interface{}
	err = json.Unmarshal([]byte(data.EventData), &d)
	if err != nil {
		return
	}
	taskData := task.Inputs(d)
	var taskDataJSON []byte
	taskDataJSON, _ = json.Marshal(taskData)
	log.Println("Trigger task", task.Name)
	client.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: task.Service,
		TaskKey:   task.Name,
		TaskData:  string(taskDataJSON),
	})
	return
}
