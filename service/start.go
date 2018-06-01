package service

import (
	"errors"
	"strings"
	"sync"
	"time"

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
	var mutex sync.Mutex
	i := 0
	var wg sync.WaitGroup
	for name, dependency := range service.GetDependencies() {
		d := dependencyDetails{
			namespace:      service.namespace(),
			dependencyName: name,
			serviceName:    service.Name,
		}
		wg.Add(1)
		go func(service *Service, d dependencyDetails, name string, i int) {
			serviceID, errStart := dependency.Start(service, d, networkID)
			mutex.Lock()
			serviceIDs[i] = serviceID
			if errStart != nil && err == nil {
				err = errStart
			}
			mutex.Unlock()
			wg.Done()
		}(service, d, name, i)
		i++
	}
	wg.Wait()
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
	namespace := []string{details.namespace, details.dependencyName} //TODO: refacto namespace
	serviceID, err = container.StartService(container.ServiceOptions{
		Namespace: namespace,
		Labels: map[string]string{
			"mesg.service": details.serviceName,
		},
		Image: dependency.Image,
		Args:  strings.Fields(dependency.Command),
		Env: []string{
			"MESG_ENDPOINT=" + viper.GetString(config.APIServiceTargetSocket),
			"MESG_ENDPOINT_TCP=mesg-daemon:50052",
		},
		Mounts: append(extractVolumes(service, dependency, details), container.Mount{
			Source: viper.GetString(config.APIServiceSocketPath),
			Target: viper.GetString(config.APIServiceTargetPath),
		}),
		Ports:      extractPorts(dependency),
		NetworksID: []string{networkID, sharedNetworkID},
	})
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(namespace, container.RUNNING, time.Minute) //TODO: be careful with timeout
	return
}
