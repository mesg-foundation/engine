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

// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStopRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	service.Start()
	err := service.Stop()
	require.NoError(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopNonRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopNonRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	err := service.Stop()
	require.NoError(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopDependency(t *testing.T) {
	c := newIntegrationContainer(t)
	service, _ := FromService(&Service{
		Name: "TestStopDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	networkID, err := c.CreateNetwork(service.namespace())
	require.NoError(t, err)
	defer c.DeleteNetwork(service.namespace(), container.EventDestroy)
	dep := service.Dependencies[0]
	dep.Start(networkID)
	err = dep.Stop()
	require.NoError(t, err)
	status, _ := dep.Status()
	require.Equal(t, container.STOPPED, status)
}

func TestIntegrationNetworkDeleted(t *testing.T) {
	c := newIntegrationContainer(t)
	service, _ := FromService(&Service{
		Name: "TestNetworkDeleted",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	service.Start()
	service.Stop()
	n, err := c.FindNetwork(service.namespace())
	require.Empty(t, n)
	require.Error(t, err)
}
