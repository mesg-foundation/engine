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

// StopService stops service serviceID.
func (a *API) StopService(serviceID string) error {
	return newServiceStopper(a).Stop(serviceID)
}

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
