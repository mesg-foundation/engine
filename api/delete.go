package api

import (
	"github.com/mesg-foundation/core/service"
)

// DeleteService stops and deletes service serviceID.
func (a *API) DeleteService(serviceID string) error {
	return newServiceDeletor(a).Delete(serviceID)
}

// serviceDeletor provides functionalities to delete a MESG service.
type serviceDeletor struct {
	api *API
}

// newServiceDeletor creates a new serviceDeletor with given.
func newServiceDeletor(api *API) *serviceDeletor {
	return &serviceDeletor{
		api: api,
	}
}

// Delete stops and deletes service serviceID.
func (d *serviceDeletor) Delete(serviceID string) error {
	s, err := d.api.db.Get(serviceID)
	if err != nil {
		return err
	}
	s, err = service.FromService(s, service.ContainerOption(d.api.container))
	if err != nil {
		return err
	}
	if err := s.Stop(); err != nil {
		return err
	}
	return d.api.db.Delete(serviceID)
}
