package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
)

// serviceLister provides functionalities to list MESG services.
type serviceLister struct {
	api *API
}

// newServiceLister creates a new serviceLister with given api and filters.
func newServiceLister(api *API, filters ...ListServicesFilter) *serviceLister {
	l := &serviceLister{
		api: api,
	}
	for _, filter := range filters {
		filter(l)
	}
	return l
}

// Lists services.
func (l *serviceLister) List() ([]*service.Service, error) {
	ss, err := services.All()
	if err != nil {
		return nil, err
	}
	var services []*service.Service
	for _, s := range ss {
		s, err = service.FromService(s, service.ContainerOption(l.api.container))
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
