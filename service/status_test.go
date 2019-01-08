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
	"errors"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestUnknownServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		statusErr     = errors.New("ops")
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestUnknownServiceStatus",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.UNKNOWN, statusErr)

	status, err := s.Status()
	require.Equal(t, statusErr, err)
	require.Equal(t, UNKNOWN, status)

	mc.AssertExpectations(t)
}

func TestStoppedServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStoppedServiceStatus",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.STOPPED, nil)

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, STOPPED, status)

	mc.AssertExpectations(t)
}

func TestRunningServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestRunningServiceStatus",
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

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, RUNNING, status)

	mc.AssertExpectations(t)
}

func TestPartialServiceStatus(t *testing.T) {
	var (
		dependencyKey  = "1"
		dependencyKey2 = "2"
		s, mc          = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestPartialServiceStatus",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
				{
					Key:   dependencyKey2,
					Image: "http-server",
				},
			},
		})
	)

	var (
		d, _  = s.getDependency(dependencyKey)
		d2, _ = s.getDependency(dependencyKey2)
	)

	mc.On("Status", d.namespace()).Once().Return(container.RUNNING, nil)
	mc.On("Status", d2.namespace()).Once().Return(container.STOPPED, nil)

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, PARTIAL, status)

	mc.AssertExpectations(t)
}

func TestDependencyStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestDependencyStatus",
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

	status, err := d.Status()
	require.NoError(t, err)
	require.Equal(t, container.RUNNING, status)

	mc.AssertExpectations(t)
}
