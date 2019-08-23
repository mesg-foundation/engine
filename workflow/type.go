package workflow

import "github.com/mesg-foundation/engine/hash"

// Workflow describes a workflow of a service
type Workflow struct {
	Graph
	Hash hash.Hash `hash:"-" validate:"required"`
	Key  string    `hash:"name:1" validate:"required"`
}

// Task is a type of node that triggers an execution
type Task struct {
	Key          string    `hash:"name:1" validate:"required"`
	InstanceHash hash.Hash `hash:"name:2" validate:"required"`
	TaskKey      string    `hash:"name:3" validate:"required,printascii"`
}

// Result is a type of node that listen for an result
type Result struct {
	InstanceHash hash.Hash `hash:"name:1" validate:"required"`
	TaskKey      string    `hash:"name:2" validate:"printascii,required_without=EventKey"`
}

// Event is a type of node that listen for an event
type Event struct {
	InstanceHash hash.Hash `hash:"name:1" validate:"required"`
	EventKey     string    `hash:"name:3" validate:"printascii,required_without=TaskKey"`
}

// Mapping is a type of Node that transform data
type Mapping struct {
	Inputs []struct {
		Key string `hash:"name:1" validate:"required"`
		Ref struct {
			NodeKey string `hash:"name:1" validate:"required"`
			Key     string `hash:"name:2" validate:"required"`
		} `hash:"name:2" validate:"required"`
	} `hash:"name:1" validate:"dive,required"`
}
