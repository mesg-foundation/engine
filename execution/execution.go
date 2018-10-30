package execution

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

// Status stores the state of an execution
type Status int

// Status for an execution
// Created 		=> The execution is created but not yet processed
// InProgress => The execution is being processed
// Completed	=> The execution is completed
const (
	Created Status = iota
	InProgress
	Completed
)

// Execution stores all informations about executions.
type Execution struct {
	ID                string                 `hash:"-"`
	Status            Status                 `hash:"-"`
	Service           *service.Service       `hash:"service"`
	TaskKey           string                 `hash:"taskKey"`
	Tags              []string               `hash:"tags"`
	Inputs            map[string]interface{} `hash:"inputs"`
	Output            string                 `hash:"-"`
	OutputData        map[string]interface{} `hash:"-"`
	CreatedAt         time.Time              `hash:"createdAt"`
	ExecutedAt        time.Time              `hash:"-"`
	ExecutionDuration time.Duration          `hash:"-"`
}

// New a record in the database to store this execution and returns the id
// returns an error if any problem happen with the database
// returns an error if inputs are invalid
func New(service *service.Service, taskKey string, inputs map[string]interface{}, tags []string) (*Execution, error) {
	task, err := service.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	if err := task.RequireInputs(inputs); err != nil {
		return nil, err
	}
	return &Execution{
		Service:   service,
		Inputs:    inputs,
		TaskKey:   taskKey,
		Tags:      tags,
		CreatedAt: time.Now(),
		Status:    Created,
	}, err
}

// Execute a given execution
// Returns an error if the execution doesn't exists in the database
// Returns an error if the status of the execution is different of `Created`
func (execution *Execution) Execute() error {
	if execution.Status != Created {
		return StatusError{
			ActualStatus:   execution.Status,
			ExpectedStatus: Created,
		}
	}
	execution.ExecutedAt = time.Now()
	execution.Status = InProgress
	return nil
}

// Complete verifies the output associated to the execution and save this to the database
// Returns an error if the executionID doesn't exists
// Returns an error if the execution is not `InProgress`
// Returns an error if the `outputKey` or `outputData` are not valid
func (execution *Execution) Complete(outputKey string, outputData map[string]interface{}) error {
	if execution.Status != InProgress {
		return StatusError{
			ActualStatus:   execution.Status,
			ExpectedStatus: InProgress,
		}
	}
	task, err := execution.Service.GetTask(execution.TaskKey)
	if err != nil {
		return err
	}
	output, err := task.GetOutput(outputKey)
	if err != nil {
		return err
	}
	if err := output.RequireData(outputData); err != nil {
		return err
	}

	execution.ExecutionDuration = time.Since(execution.ExecutedAt)
	execution.Output = outputKey
	execution.OutputData = outputData
	execution.Status = Completed

	return nil
}
