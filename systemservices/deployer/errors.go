package deployer

import (
	"fmt"
)

// systemServiceNotFoundError is returned when an expected
// system service is not found.
type systemServiceNotFoundError struct {
	name string
}

func (e *systemServiceNotFoundError) Error() string {
	return fmt.Sprintf("system service %q not found", e.name)
}
