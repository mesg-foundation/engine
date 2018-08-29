package execution

import (
	"encoding/json"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/hash"
)

// Create creates an execution with a unique ID and puts it in the pending list.
func Create(s *service.Service, taskKey string, taskInputs map[string]interface{}, tags []string) (*Execution, error) {
	task, ok := s.Tasks[taskKey]
	if !ok {
		return nil, &service.TaskNotFoundError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
		}
	}
	warnings := s.ValidateParametersSchema(task.Inputs, taskInputs)
	if len(warnings) > 0 {
		return nil, &service.InvalidTaskInputError{
			TaskKey:  taskKey,
			Warnings: warnings,
		}
	}
	execution := &Execution{
		Service:   s,
		Inputs:    taskInputs,
		Task:      taskKey,
		Tags:      tags,
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
