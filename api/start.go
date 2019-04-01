package api

import (
	"github.com/mesg-foundation/core/service"
)

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	sr, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	sr, err = service.FromService(sr, service.ContainerOption(a.container))
	if err != nil {
		return err
	}
	_, err = sr.Start()
	return err
}
