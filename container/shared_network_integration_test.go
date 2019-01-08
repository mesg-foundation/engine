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
	"context"
	"testing"

	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

func removeSharedNetworkIfExist(c *DockerContainer) error {
	if _, err := c.sharedNetwork(); err != nil {
		if docker.IsErrNotFound(err) {
			return nil
		}
		return err
	}
	return c.client.NetworkRemove(context.Background(), c.Namespace([]string{}))
}

func TestIntegrationCreateSharedNetworkIfNeeded(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	err = removeSharedNetworkIfExist(c)
	require.NoError(t, err)
	err = c.createSharedNetworkIfNeeded()
	require.NoError(t, err)
}

func TestIntegrationSharedNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	network, err := c.sharedNetwork()
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationSharedNetworkID(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	networkID, err := c.SharedNetworkID()
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}
