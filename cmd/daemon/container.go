package daemon

import (
	"context"

	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"

	"github.com/docker/docker/api/types/swarm"
	"github.com/fsouza/go-dockerclient"
)

func container() (*docker.APIContainers, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, nil
	}
	res, err := client.ListContainers(docker.ListContainersOptions{
		Context: context.Background(),
		Limit:   1,
		Filters: map[string][]string{
			"ancestor": []string{viper.GetString(config.DaemonImage)},
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
