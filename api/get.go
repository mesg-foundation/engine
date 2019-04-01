package api

import (
	"github.com/mesg-foundation/core/service"
)

// GetService returns service serviceID.
func (a *API) GetService(serviceID string) (*service.Service, error) {
	s, err := a.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	return service.FromService(s, service.ContainerOption(a.container))
}
