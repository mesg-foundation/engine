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

package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
)

// Stop stops a service.
func (s *Service) Stop() error {
	status, err := s.Status()
	if err != nil || status == STOPPED {
		return err
	}

	if err := s.StopDependencies(); err != nil {
		return err
	}
	return s.container.DeleteNetwork(s.namespace(), container.EventDestroy)
}

// StopDependencies stops all dependencies.
func (s *Service) StopDependencies() error {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var err error
	for _, dep := range s.Dependencies {
		wg.Add(1)
		go func(d *Dependency) {
			defer wg.Done()
			errStop := d.Stop()
			mutex.Lock()
			defer mutex.Unlock()
			if errStop != nil && err == nil {
				err = errStop
			}
		}(dep)
	}
	wg.Wait()
	return err
}

// Stop stops a dependency.
func (d *Dependency) Stop() error {
	return d.service.container.StopService(d.namespace())
}
