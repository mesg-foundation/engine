package service

import (
	"errors"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start a service
func (service *Service) Start() (serviceIDs []string, err error) {
	if service.IsRunning() {
		return
	}
	// If there is one but not all services running stop to restart all
	if service.IsPartiallyRunning() {
		err = service.StopDependencies()
		if err != nil {
			return
		}
	}
	networkID, err := container.CreateNetwork([]string{service.namespace()})
	if err != nil {
		return
	}
	serviceIDs = make([]string, len(service.GetDependencies()))
	i := 0
	for name, dependency := range service.GetDependencies() {
		serviceIDs[i], err = dependency.Start(service, dependencyDetails{
			namespace:      service.namespace(),
			dependencyName: name,
			serviceName:    service.Name,
		}, networkID)
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
func (dependency *Dependency) Start(service *Service, details dependencyDetails, networkID string) (serviceID string, err error) {
	if networkID == "" {
		panic(errors.New("Network ID should never be null"))
	}
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		return
	}
	return container.StartService(swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: container.Namespace([]string{details.namespace, details.dependencyName}),
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
					"MESG_ENDPOINT_TCP=mesg-daemon:50052",
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
				Target: networkID,
			},
			swarm.NetworkAttachmentConfig{
				Target: sharedNetworkID,
			},
		},
	})
}
