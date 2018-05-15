package service

import (
	"errors"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/cmd/daemon"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"
)

// Start a service
func (service *Service) Start() (dockerServices []*swarm.Service, err error) {
	if service.IsRunning() {
		return
	}
	// If there is one but not all services running stop to restart all
	if service.IsPartiallyRunning() {
		service.Stop()
	}
	dockerServices = make([]*swarm.Service, len(service.GetDependencies()))
	i := 0
	for name, dependency := range service.GetDependencies() {
		dockerServices[i], err = dependency.Start(service, dependencyDetails{
			namespace:      service.namespace(),
			dependencyName: name,
			serviceName:    service.Name,
		})
		i++
		if err != nil {
			break
		}
	}
	// Disgrasfully close the service because there is an error
	if err != nil {
		service.Stop()
	}
	return
}

type dependencyDetails struct {
	namespace      string
	dependencyName string
	serviceName    string
}

// Start will start a dependency container
func (dependency *Dependency) Start(service *Service, details dependencyDetails) (dockerService *swarm.Service, err error) {
	cli, err := dockerCli()
	if err != nil {
		return
	}
	daemonContainer, err := daemon.Container()
	if err != nil {
		return
	}
	networkContainer := daemonContainer.Networks.Networks["mesg-shared-network"]
	if networkContainer.IPAddress == "" {
		err = errors.New("Network 'mesg-shared-network' not found")
		return
	}

	return cli.CreateService(docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: strings.Join([]string{details.namespace, details.dependencyName}, "_"),
				Labels: map[string]string{
					"com.docker.stack.image":     dependency.Image,
					"com.docker.stack.namespace": details.namespace,
					"mesg.service":               details.serviceName,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: dependency.Image,
					Args:  strings.Fields(dependency.Command),
					Env: []string{
						"MESG_ENDPOINT=" + viper.GetString(config.APIServiceTargetSocket),
						"MESG_ENDPOINT_TCP=" + networkContainer.IPAddress + "" + viper.GetString(config.APIClientTarget),
					},
					Mounts: append(extractVolumes(service, dependency, details), mount.Mount{
						Source: viper.GetString(config.APIServiceSocketPath),
						Target: viper.GetString(config.APIServiceTargetPath),
					}),
					Labels: map[string]string{
						"com.docker.stack.namespace": details.namespace,
					},
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: extractPorts(dependency),
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{
					Target: "mesg-shared-network", // TODO: to set a the same constant as daemon package
				},
			},
		},
	})
}
