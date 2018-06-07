package daemon

import (
	"github.com/mesg-foundation/core/container"
)

// IsRunning returns true if the MESG Core docker service is running
func IsRunning() (running bool, err error) {
	return container.IsServiceRunning(Namespace())
}

// IsStopped returns true if the MESG Core docker service is stopped
func IsStopped() (stopped bool, err error) {
	return container.IsServiceStopped(Namespace())
}
