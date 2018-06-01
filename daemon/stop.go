package daemon

import (
	"time"

	"github.com/mesg-foundation/core/container"
)

// Stop the daemon docker
func Stop() (err error) {
	stopped, err := IsStopped()
	if err != nil || stopped == true {
		return
	}
	err = container.StopService(Namespace())
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(Namespace(), container.STOPPED, time.Minute)
	return
}
