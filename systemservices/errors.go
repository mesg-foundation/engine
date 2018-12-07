package systemservices

import (
	"fmt"
)

// SystemServiceNotFoundError is returned when an expected
// system service is not found.
type SystemServiceNotFoundError struct {
	Name string
}

func (e *SystemServiceNotFoundError) Error() string {
	return fmt.Sprintf("System service %q not found", e.Name)
}
