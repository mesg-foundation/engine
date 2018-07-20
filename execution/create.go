package execution

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/hash"
)

// Create an execution with a unique ID and put it in the pending list
func Create(serviceForExecution *service.Service, task string, inputs map[string]interface{}) (execution *Execution, err error) {
	serviceTask, taskFound := serviceForExecution.Tasks[task]
	if !taskFound {
		err = errors.New("Task " + task + " doesn't exists in service " + serviceForExecution.Name)
		return
	}
	if !serviceTask.IsValid(inputs) {
		errorString := "Invalid inputs: "
		for _, warning := range serviceTask.Validate(inputs) {
			errorString = errorString + " " + warning.String()
		}
		return nil, errors.New(errorString)
	}
	execution = &Execution{
		Service:   serviceForExecution,
		Inputs:    inputs,
		Task:      task,
		CreatedAt: time.Now(),
	}
	execution.ID, err = generateID(execution)
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	pendingExecutions[execution.ID] = execution
	log.Println("[PENDING]", task)
	return
}

func generateID(execution *Execution) (id string, err error) {
	inputs, err := json.Marshal(execution.Inputs)
	if err != nil {
		return
	}
	id = hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
		string(inputs),
	})
	return
}
