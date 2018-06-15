package daemon

import (
	"github.com/mesg-foundation/core/container"
)

// Stop the MESG Core docker container
func Stop() (err error) {
	status, err := Status()
	if err != nil || status == container.STOPPED {
		return
	}
	err = container.StopService(Namespace())
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(Namespace(), container.STOPPED)
	return
}
