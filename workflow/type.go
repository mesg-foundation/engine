package workflow

import "github.com/mesg-foundation/engine/hash"

// These structs are temporary and will be part of the service definition
// TODO: Move to service struct

type workflow struct {
	Trigger trigger
	Task    task
}

type task struct {
	InstanceHash hash.Hash
	TaskKey      string
}

// Trigger is an event that triggers a workflow
type trigger struct {
	InstanceHash hash.Hash
	EventKey     string
}
