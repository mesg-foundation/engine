package workflow

import (
	"github.com/mesg-foundation/core/x/xstrings"
)

// Workflow represents a new workflow.
type Workflow struct {
	// Hash of the workflow.
	Hash string

	// Key of the workflow.
	Key string

	// Trigger of the workflow.
	Trigger Trigger

	// Tasks of the workflow.
	Tasks []Task
}

// Match checks if the workflow's trigger matches with given trigger spec.
func (w *Workflow) Match(trigger Trigger) bool {
	return xstrings.SliceContains([]string{"*", "", trigger.InstanceHash}, w.Trigger.InstanceHash) &&
		xstrings.SliceContains([]string{"*", "", trigger.EventKey}, w.Trigger.EventKey) &&
		xstrings.SliceContains([]string{"*", "", trigger.Filter.TaskKey}, w.Trigger.Filter.TaskKey)
}

// Task to execute.
type Task struct {
	InstanceHash string
	Key          string
}

// Trigger is a pattern to match first in order to start workflow's lifecycle.
type Trigger struct {
	InstanceHash string
	EventKey     string
	Filter       Filter
}

// Filter keeps optional trigger conditions.
type Filter struct {
	TaskKey string
}
