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

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	return newServiceStarter(a).Start(serviceID)
}

// serviceStarter provides functionalities to start a MESG service.
type serviceStarter struct {
	api *API
}

// newServiceStarter creates a new serviceStarter with given api.
func newServiceStarter(api *API) *serviceStarter {
	return &serviceStarter{
		api: api,
	}
}

// Start starts service serviceID.
func (s *serviceStarter) Start(serviceID string) error {
	sr, err := s.api.db.Get(serviceID)
	if err != nil {
		return err
	}
	sr, err = service.FromService(sr, service.ContainerOption(s.api.container))
	if err != nil {
		return err
	}
	_, err = sr.Start()
	return err
}
