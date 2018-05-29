package daemon

import (
	"bytes"

	"github.com/mesg-foundation/core/container"
)

// Logs return the daemon's docker service logs
func Logs(stream *bytes.Buffer) (err error) {
	return container.ServiceLogs(Namespace(), stream)
}
