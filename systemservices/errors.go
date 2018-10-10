package systemservices

import (
	"fmt"
)

// notDirectoryError is returned when a dir expected but
// instead a file found.
type notDirectoryError struct {
	fileName string
}

func (e *notDirectoryError) Error() string {
	return fmt.Sprintf("%q not a directory", e.fileName)
}

// systemServiceNotFound returned when a system service is not found.
type systemServiceNotFound struct {
	name string
}

func (e *systemServiceNotFound) Error() string {
	return fmt.Sprintf("system service %q not found", e.name)
}
