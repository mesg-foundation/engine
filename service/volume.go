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
	"github.com/mesg-foundation/core/x/xerrors"
)

// DeleteVolumes deletes the data volumes of service and its dependencies.
func (s *Service) DeleteVolumes() error {
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, d := range s.Dependencies {
		wg.Add(1)
		go func(d *Dependency) {
			defer wg.Done()
			if err := d.DeleteVolumes(); err != nil {
				errs.Append(err)
			}
		}(d)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}

// DeleteVolumes deletes the data volumes of service's dependency.
func (d *Dependency) DeleteVolumes() error {
	volumes := d.extractVolumes()
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, mount := range volumes {
		wg.Add(1)
		go func(mount container.Mount) {
			defer wg.Done()
			if err := d.service.container.DeleteVolume(mount.Source); err != nil {
				errs.Append(err)
			}
		}(mount)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
