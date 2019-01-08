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

package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// SharedNetworkID returns the ID of the shared network.
func (c *DockerContainer) SharedNetworkID() (networkID string, err error) {
	network, err := c.sharedNetwork()
	if err != nil {
		return "", err
	}
	return network.ID, nil
}

func (c *DockerContainer) createSharedNetworkIfNeeded() error {
	network, err := c.sharedNetwork()
	if err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if network.ID != "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()

	// Create the new network needed to run containers.
	namespace := c.Namespace([]string{})
	_, err = c.client.NetworkCreate(ctx, namespace, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return err
}

// sharedNetwork returns the shared network created to connect services and MESG Core.
func (c *DockerContainer) sharedNetwork() (network types.NetworkResource, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.NetworkInspect(ctx, c.Namespace([]string{}), types.NetworkInspectOptions{})
}
