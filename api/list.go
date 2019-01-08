// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/mesg-foundation/core/service"
)

// ListServicesFilter is a filter func for listing services.
type ListServicesFilter func(*serviceLister)

// ListServices returns services matches with filters.
func (a *API) ListServices(filters ...ListServicesFilter) ([]*service.Service, error) {
	return newServiceLister(a, filters...).List()
}

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
	ss, err := l.api.db.All()
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
