package api

import (
	"github.com/mesg-foundation/core/database/services"
)

// serviceStopper provides functionalities to start a MESG service.
type serviceStopper struct {
	api *API
}

// newServiceStopper creates a new serviceStopper with given api.
func newServiceStopper(api *API) *serviceStopper {
	return &serviceStopper{
		api: api,
	}
}

// Stop stops service serviceID.
func (s *serviceStopper) Stop(serviceID string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	return service.Stop()
}
