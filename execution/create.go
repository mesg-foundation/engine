package execution

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/hash"
)

// Create an execution with a unique ID and put it in the pending list
func Create(serviceForExecution *service.Service, task string, inputs map[string]interface{}) (*Execution, error) {
	serviceTask, taskFound := serviceForExecution.Tasks[task]
	if !taskFound {
		return nil, &service.TaskNotFoundError{
			Service: serviceForExecution,
			TaskKey: task,
		}
	}
	if !serviceTask.IsValid(inputs) {
		return nil, &service.InvalidTaskInputError{
			Task:   serviceTask,
			Key:    task,
			Inputs: inputs,
		}
	}
	execution := &Execution{
		Service:   serviceForExecution,
		Inputs:    inputs,
		Task:      task,
		CreatedAt: time.Now(),
	}
	var err error
	execution.ID, err = generateID(execution)
	if err != nil {
		return nil, err
	}
	mu.Lock()
	defer mu.Unlock()
	pendingExecutions[execution.ID] = execution
	log.Println("[PENDING]", task)
	return execution, err
}

func generateID(execution *Execution) (string, error) {
	inputs, err := json.Marshal(execution.Inputs)
	if err != nil {
		return "", err
	}
	return hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
		string(inputs),
	}), nil
}
