package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xstructhash"
)

// Start starts the service.
func (s *Service) Start() (serviceIDs []string, err error) {
	status, err := s.Status()
	if err != nil || status == RUNNING {
		return nil, err //TODO: if the service is already running, serviceIDs should be returned.
	}
	// If there is one but not all services running stop to restart all
	if status == PARTIAL {
		if err := s.StopDependencies(); err != nil {
			return nil, err
		}
	}
	networkID, err := defaultContainer.CreateNetwork(s.namespace())
	if err != nil {
		return nil, err
	}
	var (
		mutex sync.Mutex
		wg    sync.WaitGroup
	)
	serviceIDs = make([]string, len(s.Dependencies))
	for i, dependency := range s.Dependencies {
		wg.Add(1)
		go func(dep *Dependency, i int) {
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
		s.Stop()
	}
	return serviceIDs, err
}

// Start starts a dependency container.
func (d *Dependency) Start(networkID string) (containerServiceID string, err error) {
	sharedNetworkID, err := defaultContainer.SharedNetworkID()
	if err != nil {
		return "", err
	}
	mounts, err := d.extractVolumes()
	if err != nil {
		return "", err
	}
	endpoint := "mesg-core:50052" // TODO: should get this from daemon namespace and config
	return defaultContainer.StartService(container.ServiceOptions{
		Namespace: d.namespace(),
		Labels: map[string]string{
			"mesg.service": d.service.Name,
			"mesg.hash":    d.service.ID,
		},
		Image: d.Image,
		Args:  strings.Fields(d.Command),
		Env: container.MapToEnv(map[string]string{
			"MESG_TOKEN":        d.service.ID,
			"MESG_ENDPOINT":     endpoint,
			"MESG_ENDPOINT_TCP": endpoint,
		}),
		Mounts:     mounts,
		Ports:      d.extractPorts(),
		NetworksID: []string{networkID, sharedNetworkID},
	})
}

func (d *Dependency) extractPorts() []container.Port {
	ports := make([]container.Port, len(d.Ports))
	for i, p := range d.Ports {
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
func (d *Dependency) extractVolumes() ([]container.Mount, error) {
	volumes := make([]container.Mount, 0)
	for _, volume := range d.Volumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(d.service, d.Key, volume),
			Target: volume,
		})
	}
	for _, depName := range d.VolumesFrom {
		dep, err := d.service.getDependency(depName)
		if err != nil {
			return nil, err
		}
		for _, volume := range dep.Volumes {
			volumes = append(volumes, container.Mount{
				Source: volumeKey(d.service, depName, volume),
				Target: volume,
			})
		}
	}
	return volumes, nil
}

func volumeKey(s *Service, dependency string, volume string) string {
	return xstructhash.Hash([]string{
		s.ID,
		dependency,
		volume,
	}, 1)
}
