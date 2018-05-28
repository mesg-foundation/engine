package daemon

import (
	"bytes"

	"github.com/mesg-foundation/core/docker"
)

// Logs return the daemon's docker service logs
func Logs(stream *bytes.Buffer) (err error) {
	return docker.ServiceLogs(Namespace(), stream)
}
