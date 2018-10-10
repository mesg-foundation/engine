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
