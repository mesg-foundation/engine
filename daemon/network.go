package daemon

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/docker"
)

func IP() (daemonIP string, err error) {
	daemonContainer, err := Container() // TODO: should try to use the service and then delete the func DaemonContainer()
	if err != nil {
		return
	}
	if daemonContainer == nil {
		err = errors.New("Daemon container is not found")
		return
	}
	networkContainer := daemonContainer.Networks.Networks["mesg-shared-network"]
	if networkContainer.IPAddress == "" {
		err = errors.New("Network 'mesg-shared-network' not found")
		return
	}
	daemonIP = networkContainer.IPAddress
	return
}

func Container() (*godocker.APIContainers, error) {
	client, err := docker.Client()
	if err != nil {
		return nil, nil
	}
	res, err := client.ListContainers(godocker.ListContainersOptions{
		Context: context.Background(),
		Limit:   1,
		Filters: map[string][]string{
			"ancestor": []string{image},
			"status":   []string{"running"},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func Service() (*swarm.Service, error) {
	return docker.FindService([]string{name})
}

func Network() (network *godocker.Network, err error) {
	return docker.FindNetwork(sharedNetwork)
}

func DaemonStart() {

}
