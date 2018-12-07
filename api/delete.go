package api

import (
	"github.com/mesg-foundation/core/service"
)

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// service and its dependencies also will be deleted.
func (a *API) DeleteService(serviceID string, deleteData bool) error {
	return newServiceDeletor(a).Delete(serviceID, deleteData)
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
// when deleteData is enabled, any persistent data that belongs to
// service and its dependencies also will be deleted.
func (d *serviceDeletor) Delete(serviceID string, deleteData bool) error {
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
	// delete volumes first before the service. this way if
	// deleting volumes fails, process can be retried by the user again
	// because service still will be in the db.
	if deleteData {
		if err := s.DeleteVolumes(); err != nil {
			return err
		}
	}
	return d.api.db.Delete(serviceID)
}
