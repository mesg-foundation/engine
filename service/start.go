package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xstructhash"
)

// Start starts the service.
func (service *Service) Start() (serviceIDs []string, err error) {
	status, err := service.Status()
	if err != nil || status == RUNNING {
		return nil, err //TODO: if the service is already running, serviceIDs should be returned.
	}
	// If there is one but not all services running stop to restart all
	if status == PARTIAL {
		if err := service.StopDependencies(); err != nil {
			return nil, err
		}
	}
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	if err != nil {
		return nil, err
	}
	var (
		mutex                   sync.Mutex
		wg                      sync.WaitGroup
		dependenciesFromService = service.DependenciesFromService()
	)
	serviceIDs = make([]string, len(dependenciesFromService))
	for i, dependency := range dependenciesFromService {
		wg.Add(1)
		go func(dep *DependencyFromService, i int) {
			defer wg.Done()
			serviceID, errStart := dep.Start(networkID)
			mutex.Lock()
			defer mutex.Unlock()
			serviceIDs[i] = serviceID
			if errStart != nil && err == nil {
				err = errStart
			}
		}(dependency, i)
	}
	wg.Wait()
	// Gracefully stop the service because there is an error
	if err != nil {
		service.Stop()
	}
	return serviceIDs, err
}

// Start starts a dependency container.
func (dependency *DependencyFromService) Start(networkID string) (containerServiceID string, err error) {
	if networkID == "" {
		return "", errors.New("Network ID should never be null")
	}
	service := dependency.Service
	if service == nil {
		return "", errors.New("Service is nil")
	}
	sharedNetworkID, err := defaultContainer.SharedNetworkID()
	if err != nil {
		return "", err
	}
	mounts, err := dependency.extractVolumes()
	if err != nil {
		return "", err
	}
	apiPort, err := config.APIPort().GetValue()
	if err != nil {
		return "", err
	}
	endpoint := "mesg-core:" + apiPort // TODO: should get this from daemon namespace and config
	return defaultContainer.StartService(container.ServiceOptions{
		Namespace: dependency.namespace(),
		Labels: map[string]string{
			"mesg.service": service.Name,
			"mesg.hash":    service.Hash(),
		},
		Image: dependency.Image,
		Args:  strings.Fields(dependency.Command),
		Env: container.MapToEnv(map[string]string{
			"MESG_TOKEN":        service.Hash(),
			"MESG_ENDPOINT":     endpoint,
			"MESG_ENDPOINT_TCP": endpoint,
		}),
		Mounts:     mounts,
		Ports:      dependency.extractPorts(),
		NetworksID: []string{networkID, sharedNetworkID},
	})
}

func (dependency *Dependency) extractPorts() []container.Port {
	ports := make([]container.Port, len(dependency.Ports))
	for i, p := range dependency.Ports {
		split := strings.Split(p, ":")
		from, _ := strconv.ParseUint(split[0], 10, 64)
		to := from
		if len(split) > 1 {
			to, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = container.Port{
			Target:    uint32(to),
			Published: uint32(from),
		}
	}
	return ports
}

// TODO: add test and hack for MkDir in CircleCI
func (dependency *DependencyFromService) extractVolumes() ([]container.Mount, error) {
	service := dependency.Service
	if service == nil {
		return nil, errors.New("Service is nil")
	}
	volumes := make([]container.Mount, 0)
	for _, volume := range dependency.Volumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(service, dependency.Name, volume),
			Target: volume,
		})
	}
	for _, depName := range dependency.Volumesfrom {
		dep := service.Dependencies[depName]
		if dep == nil {
			return nil, fmt.Errorf("Dependency %s do not exist", depName)
		}
		for _, volume := range dep.Volumes {
			volumes = append(volumes, container.Mount{
				Source: volumeKey(service, depName, volume),
				Target: volume,
			})
		}
	}
	return volumes, nil
}

func volumeKey(s *Service, dependency string, volume string) string {
	return xstructhash.Hash([]string{
		s.Hash(),
		dependency,
		volume,
	}, 1)
}
