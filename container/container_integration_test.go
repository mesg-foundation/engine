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

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationFindContainerNotExisting(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	_, err = c.FindContainer([]string{"TestFindContainerNotExisting"})
	require.Error(t, err)
}

func TestIntegrationFindContainer(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestFindContainer"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	container, err := c.FindContainer(namespace)
	require.NoError(t, err)
	require.NotEqual(t, "", container.ID)
}

func TestIntegrationFindContainerStopped(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindContainer(namespace)
	require.Error(t, err)
}

func TestIntegrationContainerStatusNeverStarted(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, STOPPED)
}

func TestIntegrationContainerStatusRunning(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, RUNNING)
}

func TestIntegrationContainerStatusStopped(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	c.waitForStatus(namespace, STOPPED)
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, STOPPED)
}
