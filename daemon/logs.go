package daemon

import (
	"io"
)

// Logs returns the core's docker service logs.
func Logs() (io.ReadCloser, error) {
	return defaultContainer.ServiceLogs(Namespace())
}
