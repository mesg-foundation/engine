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
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace()).Once().Return(nil)
	mc.On("DeleteNetwork", s.namespace(), container.EventDestroy).Once().Return(nil)

	err := s.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}

func TestStopDependency(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)
	mc.On("StopService", d.namespace()).Once().Return(nil)
	require.NoError(t, d.Stop())
	mc.AssertExpectations(t)
}
