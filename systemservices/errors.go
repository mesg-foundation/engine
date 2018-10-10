package systemservices

import (
	"fmt"
)

type notDirectoryError struct {
	fileName string
}

func (e *notDirectoryError) Error() string {
	return fmt.Sprintf("%q not a directory", e.fileName)
}
