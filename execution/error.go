package execution

import (
	"strings"

	"github.com/mesg-foundation/core/service"
)

// MissingOutputError is an error when a service doesn't contains a specific output
type MissingOutputError struct {
	Service *service.Service
	Output  string
}

func (e *MissingOutputError) Error() string {
	return strings.Join([]string{
		"Output",
		e.Output,
		"doesn't exists in service",
		e.Service.Name,
	}, " ")
}

// InvalidOutputError is an error when the outputs for one task result are not valid
type InvalidOutputError struct {
	Service  *service.Service
	Warnings []*service.ParameterWarning
}

func (e *InvalidOutputError) Error() string {
	errorString := "Invalid result: "
	for _, warning := range e.Warnings {
		errorString = errorString + " " + warning.String()
	}
	return errorString
}

// NotInQueueError is an error when trying to access an execution that doesn't exists
type NotInQueueError struct {
	ID    string
	Queue string
}

func (e *NotInQueueError) Error() string {
	return strings.Join([]string{
		"Execution",
		e.ID,
		"not in",
		e.Queue,
		"queue",
	}, " ")
}
