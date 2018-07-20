package service

import (
	"strings"
)

// MissingExecutionError is an error when an execution doesn't exists
type MissingExecutionError struct {
	ID string
}

func (e *MissingExecutionError) Error() string {
	return strings.Join([]string{
		"Execution",
		e.ID,
		"doesn't exists",
	}, " ")
}
