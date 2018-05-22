package service

import (
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/docker"
	dockerDependency "github.com/mesg-foundation/core/docker/dependency"
	dockerService "github.com/mesg-foundation/core/docker/service"
	"github.com/spf13/viper"
)

// Start a service
func (service *Service) Start(daemonIP string, sharedNetwork string) (dockerServices []*swarm.Service, err error) {
	if dockerService.IsRunning(service) {
		return
	}
	// If there is one but not all services running stop to restart all
	if dockerService.IsPartiallyRunning(service) {
		service.Stop()
	}
	dockerServices = make([]*swarm.Service, len(service.GetDependencies()))
	i := 0
	for name, dependency := range service.GetDependencies() {
		dockerServices[i], err = dependency.Start(service, dependencyDetails{
			namespace:      service.Namespace(),
			dependencyName: name,
			serviceName:    service.Name,
		}, daemonIP, sharedNetwork)
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
func (dependency *Dependency) Start(service *Service, details dependencyDetails, daemonIP string, sharedNetwork string) (dockerService *swarm.Service, err error) {
	client, err := docker.Client()
	if err != nil {
		return
	}
	return client.CreateService(godocker.CreateServiceOptions{
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
						"MESG_ENDPOINT_TCP=" + daemonIP + "" + viper.GetString(config.APIClientTarget),
					},
					Mounts: append(dockerDependency.Volumes(service, dependency, details.dependencyName), mount.Mount{
						Source: viper.GetString(config.APIServiceSocketPath),
						Target: viper.GetString(config.APIServiceTargetPath),
					}),
					Labels: map[string]string{
						"com.docker.stack.namespace": details.namespace,
					},
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: dockerDependency.Ports(dependency),
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{
					Target: sharedNetwork,
				},
			},
		},
	})
}
