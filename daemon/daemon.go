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

package daemon

import (
	"io"
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/mesg-foundation/core/x/xos"
)

// Daemon is an interface that start, stop etc core as daemon.
type Daemon interface {
	Start() error
	Stop() error
	Status() (container.StatusType, error)
	Logs() (io.ReadCloser, error)
}

// ContainerDaemon run core as container.
type ContainerDaemon struct {
	c   container.Container
	cfg *config.Config
}

// NewContainerDaemon creates new dameon that will be run in container.
func NewContainerDaemon(cfg *config.Config, c container.Container) *ContainerDaemon {
	return &ContainerDaemon{
		c:   c,
		cfg: cfg,
	}
}

// Start starts the docker core.
func (d *ContainerDaemon) Start() error {
	sharedNetworkID, err := d.c.SharedNetworkID()
	if err != nil {
		return err
	}
	_, err = d.c.StartService(d.buildServiceOptions(sharedNetworkID))
	return err
}

// Stop stops the MESG Core docker container.
func (d *ContainerDaemon) Stop() error {
	return d.c.StopService([]string{})
}

// Status returns the Status of the docker service of the daemon.
func (d *ContainerDaemon) Status() (container.StatusType, error) {
	return d.c.Status([]string{})
}

// Logs returns the core's docker service logs.
func (d *ContainerDaemon) Logs() (io.ReadCloser, error) {
	return d.c.ServiceLogs([]string{})
}

func (d *ContainerDaemon) buildServiceOptions(sharedNetworkID string) container.ServiceOptions {
	_, port, _ := xnet.SplitHostPort(d.cfg.Server.Address)
	return container.ServiceOptions{
		Namespace: []string{},
		Image:     d.cfg.Core.Image,
		Env:       xos.EnvMapToSlice(d.cfg.DaemonEnv()),
		Mounts: []container.Mount{
			{
				Source: d.cfg.Docker.Socket,
				Target: d.cfg.Docker.Socket,
				Bind:   true,
			},
			{
				Source: d.cfg.Core.Path,
				Target: d.cfg.Docker.Core.Path,
				Bind:   true,
			},
			{
				Source: filepath.Join(d.cfg.Core.Path, d.cfg.SystemServices.RelativePath),
				Target: filepath.Join(d.cfg.Docker.Core.Path, d.cfg.SystemServices.RelativePath),
				Bind:   true,
			},
		},
		Ports: []container.Port{
			{
				Target:    uint32(port),
				Published: uint32(port),
			},
		},
		Networks: []container.Network{
			{ID: sharedNetworkID},
		},
	}
}
