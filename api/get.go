package api

import "github.com/mesg-foundation/core/service"

// GetService returns service serviceID.
func (a *API) GetService(serviceID string) (*service.Service, error) {
	return newServiceGetter(a).Get(serviceID)
}
