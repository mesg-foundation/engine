package application

import "context"

// Execution is a task exection.
type Execution struct {
	// ID is execution id of task.
	ID string

	// Err filled if an error occurs during task execution.
	Err error
}

// Stream is a task execution stream.
type Stream struct {
	// Executions filled with task executions.
	Executions chan *Execution

	// Err filled when stream fails to continue.
	Err chan error

	cancel context.CancelFunc
}

// Close gracefully stops listening for events and shutdowns stream.
func (s *Stream) Close() error {
	s.cancel()
	return nil
}
