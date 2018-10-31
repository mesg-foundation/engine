package execution

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/cnf/structhash"
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
	EventID           string                 `hash:"eventID"`
	Status            Status                 `hash:"-"`
	Service           *service.Service       `hash:"service"`
	TaskKey           string                 `hash:"taskKey"`
	Tags              []string               `hash:"tags"`
	Inputs            map[string]interface{} `hash:"inputs"`
	Output            string                 `hash:"-"`
	OutputData        map[string]interface{} `hash:"-"`
	CreatedAt         time.Time              `hash:"-"`
	ExecutedAt        time.Time              `hash:"-"`
	ExecutionDuration time.Duration          `hash:"-"`
}

// New returns a new execution. It returns an error if inputs are invalid.
func New(service *service.Service, eventID string, taskKey string, inputs map[string]interface{}, tags []string) (*Execution, error) {
	task, err := service.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	if err := task.RequireInputs(inputs); err != nil {
		return nil, err
	}
	exec := &Execution{
		EventID:   eventID,
		Service:   service,
		Inputs:    inputs,
		TaskKey:   taskKey,
		Tags:      tags,
		CreatedAt: time.Now(),
		Status:    Created,
	}
	exec.ID = fmt.Sprintf("%x", sha1.Sum(structhash.Dump(exec, 1)))
	return exec, nil
}

// Execute changes executions status to in progres and update its execute time.
// It returns an error if the status is different then Created.
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

// Complete changes execution status to completed. It verifies the output.
// It returns an error if the status is different then InProgress or verification fails.
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
