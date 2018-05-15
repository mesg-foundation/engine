package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/fsouza/go-dockerclient"
)

func Container() (*docker.APIContainers, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, nil
	}
	res, err := client.ListContainers(docker.ListContainersOptions{
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

func service() (*swarm.Service, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, nil
	}
	res, err := client.ListServices(docker.ListServicesOptions{
		Context: context.Background(),
		Filters: map[string][]string{
			"name": []string{name},
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

func network() (network *docker.Network, err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return
	}

	networks, err := client.FilteredListNetworks(docker.NetworkFilterOpts{
		"name": map[string]bool{
			sharedNetwork: true,
		},
	})
	if err != nil {
		return
	}
	if len(networks) == 1 {
		network = &networks[0]
		return
	}

	return

}
