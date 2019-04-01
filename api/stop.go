package api

import (
	"github.com/mesg-foundation/core/service"
)

// StopService stops service serviceID.
func (a *API) StopService(serviceID string) error {
	sr, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	sr, err = service.FromService(sr, service.ContainerOption(a.container))
	if err != nil {
		return err
	}
	return sr.Stop()
}
