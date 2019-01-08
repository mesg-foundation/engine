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
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/systemservices"
)

// API exposes all functionalities of MESG core.
type API struct {
	db             database.ServiceDB
	execDB         database.ExecutionDB
	systemservices *systemservices.SystemServices
	container      container.Container
}

// Option is a configuration func for MESG.
type Option func(*API)

// New creates a new API with given options.
func New(db database.ServiceDB, execDB database.ExecutionDB, systemservices *systemservices.SystemServices, options ...Option) (*API, error) {
	a := &API{db: db, execDB: execDB, systemservices: systemservices}
	for _, option := range options {
		option(a)
	}
	if a.container == nil {
		var err error
		a.container, err = container.New()
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// ContainerOption configures underlying container access API.
func ContainerOption(container container.Container) Option {
	return func(a *API) {
		a.container = container
	}
}
