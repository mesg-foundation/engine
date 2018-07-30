package daemon

import (
	"io"

	"github.com/mesg-foundation/core/container"
)

// Logs return the core's docker service logs
func Logs() (io.ReadCloser, error) {
	return container.ServiceLogs(Namespace())
}
