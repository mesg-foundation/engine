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
	ServiceID         string                 `hash:"serviceID"`
	Task              *service.Task          `hash:"task"`
	Tags              []string               `hash:"tags"`
	Inputs            map[string]interface{} `hash:"inputs"`
	Output            string                 `hash:"-"`
	OutputData        map[string]interface{} `hash:"-"`
	CreatedAt         time.Time              `hash:"createdAt"`
	ExecutedAt        time.Time              `hash:"-"`
	ExecutionDuration time.Duration          `hash:"-"`
}

// DB exposes all the functionalities
type DB interface {
	Create(task *service.Task, taskInputs map[string]interface{}, tags []string) (*Execution, error)
	Find(executionID string) (*Execution, error)
	Execute(executionID string) (*Execution, error)
	Complete(executionID string, outputKey string, outputData map[string]interface{}) (*Execution, error)
	Close() error
}
