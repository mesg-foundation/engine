package daemon

import (
	"github.com/mesg-foundation/core/container"
)

// Status returns the Status of the docker service of the daemon.
func Status() (container.StatusType, error) {
	//TODO: should it be containerStatus?
	serviceStatus, err := defaultContainer.ServiceStatus(Namespace())
	if err != nil {
		return container.STOPPED, err
	}

	containerStatus, err := defaultContainer.Status(Namespace())
	if err != nil {
		return container.STOPPED, err
	}

	if serviceStatus == containerStatus {
		return serviceStatus, nil
	}
	return container.STOPPED, nil
}
