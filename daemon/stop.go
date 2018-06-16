package daemon

import (
	"github.com/mesg-foundation/core/container"
)

// Stop the MESG Core docker container
func Stop() (err error) {
	stopped, err := IsStopped()
	if err != nil || stopped == true {
		return
	}
	err = container.StopService(Namespace())
	return
}
