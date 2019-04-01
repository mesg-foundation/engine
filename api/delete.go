package api

import (
	"github.com/mesg-foundation/core/service"
)

// DeleteService stops and deletes service serviceID.
// when deleteData is enabled, any persistent data that belongs to
// the service and to its dependencies also will be deleted.
func (a *API) DeleteService(serviceID string, deleteData bool) error {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return err
	}
	s, err = service.FromService(s, service.ContainerOption(a.container))
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
	return a.db.Delete(serviceID)
}
