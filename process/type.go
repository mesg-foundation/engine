package process

import (
	"github.com/mesg-foundation/engine/filter"
	"github.com/mesg-foundation/engine/hash"
)

// Process is a graph that defines a process between multiple service instances
type Process struct {
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
	Key          string    `hash:"name:1" validate:"required"`
	InstanceHash hash.Hash `hash:"name:2" validate:"required"`
	TaskKey      string    `hash:"name:3" validate:"printascii,required"`
}

// Event is a type of node that listen for an event
type Event struct {
	Key          string    `hash:"name:1" validate:"required"`
	InstanceHash hash.Hash `hash:"name:2" validate:"required"`
	EventKey     string    `hash:"name:3" validate:"printascii,required"`
}

// Map is a type of Node that transform data
type Map struct {
	Key     string   `hash:"name:1" validate:"required"`
	Outputs []Output `hash:"name:2" validate:"dive,required"`
}

// Output as defined in a map node
type Output struct {
	Key string           `hash:"name:1" validate:"required"`
	Ref *OutputReference `hash:"name:2" validate:"required"`
}

// OutputReference of a output value to define an output
type OutputReference struct {
	NodeKey string `hash:"name:1" validate:"required"`
	Key     string `hash:"name:2" validate:"required"`
}

// Filter contains a list of conditions to apply
type Filter struct {
	filter.Filter `hash:"name:1" validate:"required"`
	Key           string `hash:"name:2" validate:"required"`
}
