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

// GetService returns service serviceID.
func (a *API) GetService(serviceID string) (*service.Service, error) {
	return newServiceGetter(a).Get(serviceID)
}

// serviceGetter provides functionalities to get a MESG service.
type serviceGetter struct {
	api *API
}

// newServiceGetter creates a new serviceGetter with given api.
func newServiceGetter(api *API) *serviceGetter {
	return &serviceGetter{
		api: api,
	}
}

// Get returns service serviceID.
func (g *serviceGetter) Get(serviceID string) (*service.Service, error) {
	s, err := g.api.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	return service.FromService(s, service.ContainerOption(g.api.container))
}
