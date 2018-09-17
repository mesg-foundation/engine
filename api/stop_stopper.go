package api

import (
	"github.com/mesg-foundation/core/service"
)

// serviceStopper provides functionalities to stop a MESG service.
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
	sr, err := s.api.db.Get(serviceID)
	if err != nil {
		return err
	}
	sr, err = service.FromService(sr, service.ContainerOption(s.api.container))
	if err != nil {
		return err
	}
	return sr.Stop()
}
