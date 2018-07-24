package service

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start a service
func (service *Service) Start() (serviceIDs []string, err error) {
	status, err := service.Status()
	if err != nil || status == RUNNING {
		return //TODO: if the service is already running, serviceIDs should be returned.
	}
	// If there is one but not all services running stop to restart all
	if status == PARTIAL {
		err = service.StopDependencies()
		if err != nil {
			return
		}
	}
	networkID, err := container.CreateNetwork(service.namespace())
	if err != nil {
		return
	}
	dependenciesFromService := service.DependenciesFromService()
	serviceIDs = make([]string, len(dependenciesFromService))
	var mutex sync.Mutex
	var wg sync.WaitGroup
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
	// Grasfully stop the service because there is an error
	if err != nil {
		service.Stop()
	}
	return
}

// Start will start a dependency container
func (dependency *DependencyFromService) Start(networkID string) (containerServiceID string, err error) {
	if networkID == "" {
		err = errors.New("Network ID should never be null")
		return
	}
	service := dependency.Service
	if service == nil {
		err = errors.New("Service is nil")
		return
	}
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		return
	}
	mounts, err := dependency.extractVolumes()
	if err != nil {
		return
	}
	containerServiceID, err = container.StartService(container.ServiceOptions{
		Namespace: dependency.namespace(),
		Labels: map[string]string{
			"mesg.service": service.Name,
			"mesg.hash":    service.Hash(),
		},
		Image: dependency.Image,
		Args:  strings.Fields(dependency.Command),
		Env: container.MapToEnv(map[string]string{
			"MESG_TOKEN":        service.Hash(),
			"MESG_ENDPOINT":     viper.GetString(config.APIServiceTargetSocket) + viper.GetString(config.APIServiceTargetPath),
			"MESG_ENDPOINT_TCP": "mesg-core:50052", // TODO: should get this from daemon namespace and config
		}),
		Mounts: append(mounts, container.Mount{
			Source: filepath.Join(viper.GetString(config.MESGPath), viper.GetString(config.APIServiceSocketPath)),
			Target: viper.GetString(config.APIServiceTargetPath),
		}),
		Ports:      dependency.extractPorts(),
		NetworksID: []string{networkID, sharedNetworkID},
	})
	return
}

func (dependency *Dependency) extractPorts() (ports []container.Port) {
	ports = make([]container.Port, len(dependency.Ports))
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
	return
}

// TODO: add test and hack for MkDir in CircleCI
func (dependency *DependencyFromService) extractVolumes() (volumes []container.Mount, err error) {
	service := dependency.Service
	if service == nil {
		err = errors.New("Service is nil")
		return
	}
	servicePath := strings.Join(service.namespace(), "-")
	volumes = make([]container.Mount, 0)
	for _, volume := range dependency.Volumes {
		path := filepath.Join(servicePath, dependency.Name, volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		volumes = append(volumes, container.Mount{
			Source: source,
			Target: volume,
		})
		// TODO: move mkdir in container package
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, depName := range dependency.Volumesfrom {
		dep := service.Dependencies[depName]
		if dep == nil {
			err = errors.New("Dependency " + depName + " do not exist")
			return
		}
		for _, volume := range dep.Volumes {
			path := filepath.Join(servicePath, depName, volume)
			source := filepath.Join(viper.GetString(config.ServicePathHost), path)
			volumes = append(volumes, container.Mount{
				Source: source,
				Target: volume,
			})
			// TODO: move mkdir in container package
			os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
		}
	}
	return
}
