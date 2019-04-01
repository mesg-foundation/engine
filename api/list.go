package api

import (
	"github.com/mesg-foundation/core/service"
)

// ListServices returns all services.
func (a *API) ListServices() ([]*service.Service, error) {
	ss, err := a.db.All()
	if err != nil {
		return nil, err
	}

	var services []*service.Service
	for _, s := range ss {
		s, err = service.FromService(s, service.ContainerOption(a.container))
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
