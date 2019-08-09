package workflow

import "github.com/mesg-foundation/engine/hash"

// TriggerType is the type for the possible triggers for a workflow
type TriggerType uint

// List of possible triggers for a workflow
const (
	EVENT TriggerType = iota + 1
	RESULT
)

// Predicate is the type of conditions that can be applied in a filter of a workflow trigger
type Predicate uint

// List of possible conditions for workflow's filter
const (
	EQ Predicate = iota + 1
)

// Workflow describes a workflow of a service
type Workflow struct {
	Hash    hash.Hash `hash:"-" validate:"required"`
	Key     string    `hash:"name:1" validate:"required"`
	Trigger Trigger   `hash:"name:2" validate:"required"`
	Nodes   []Node    `hash:"name:3" validate:"dive,required"`
	Edges   []Edge    `hash:"name:4" validate:"dive,required"`
}

// Node describes the instructions for the workflow to execute a task
type Node struct {
	Key          string    `hash:"name:1" validate:"required"`
	InstanceHash hash.Hash `hash:"name:1" validate:"required"`
	TaskKey      string    `hash:"name:2" validate:"required,printascii"`
}

// Edge describes the instructions for the workflow to execute a task
type Edge struct {
	Src string `hash:"name:1"`
	Dst string `hash:"name:2"`
}

// Trigger is an event that triggers a workflow
type Trigger struct {
	InstanceHash hash.Hash        `hash:"name:1" validate:"required"`
	Key          string           `hash:"name:2" validate:"required,printascii"`
	Type         TriggerType      `hash:"name:3" validate:"required"`
	Filters      []*TriggerFilter `hash:"name:4" validate:"dive,required"`
	InitialNode  string           `hash:"name:5"`
}

// TriggerFilter is the filter definition that can be applied to a workflow trigger
type TriggerFilter struct {
	Key       string      `hash:"name:1" validate:"required,printascii"`
	Predicate Predicate   `hash:"name:2" validate:"required"`
	Value     interface{} `hash:"name:3"`
}
