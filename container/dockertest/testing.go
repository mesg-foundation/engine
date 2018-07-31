// Package dockertest provides functionalities for mocking and faking the Docker API.
package dockertest

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

type Testing struct {
	client *client
}

func New() *Testing {
	t := &Testing{
		client: newClient(),
	}
	return t
}

func (t *Testing) Client() *client {
	return t.client
}

func (t *Testing) ProvideContainer(container types.Container) {
	t.client.m.Lock()
	defer t.client.m.Unlock()
	t.client.containers = append(t.client.containers, container)
}

func (t *Testing) ProvideContainerInspect(container types.ContainerJSON) {
	t.client.containerInspect = container
}

func (t *Testing) LastContainerList() chan types.ContainerListOptions {
	return t.client.lastContainerList
}

func (t *Testing) LastContainerInspect() chan string {
	return t.client.lastContainerInspect
}

func (t *Testing) LastSwarmInit() chan swarm.InitRequest {
	return t.client.lastSwarmInit
}

func (t *Testing) LastNetworkCreate() chan NetworkCreate {
	return t.client.lastNetworkCreate
}

type SystemInfo struct {
	Info types.Info
	Err  error
}

func (t *Testing) LastInfo() types.Info {
	return <-t.client.lastInfo
}

func (t *Testing) ProvideInfo(info types.Info, err error) {
	t.client.m.Lock()
	defer t.client.m.Unlock()
	t.client.info = info
	t.client.infoErr = err
}

func (t *Testing) LastNegotiateAPIVersion() chan struct{} {
	return t.client.lastNegotiateAPIVersion
}
