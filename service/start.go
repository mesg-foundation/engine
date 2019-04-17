package service

import (
	"strconv"
	"strings"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/mesg-foundation/core/x/xstructhash"
)

// Start starts the service.
func (s *Service) Start(c container.Container) (serviceIDs []string, err error) {
	status, err := s.Status(c)
	if err != nil || status == RUNNING {
		return nil, err //TODO: if the service is already running, serviceIDs should be returned.
	}
	// If there is one but not all services running stop to restart all
	if status == PARTIAL {
		if err := s.Stop(c); err != nil {
			return nil, err
		}
	}
	networkID, err := c.CreateNetwork(s.namespace())
	if err != nil {
		return nil, err
	}
	sharedNetworkID, err := c.SharedNetworkID()
	if err != nil {
		return nil, err
	}
	conf, err := config.Global()
	if err != nil {
		return nil, err
	}
	_, port, _ := xnet.SplitHostPort(conf.Server.Address)
	endpoint := conf.Core.Name + ":" + strconv.Itoa(port)
	// BUG: https://github.com/mesg-foundation/core/issues/382
	// After solving this by docker, switch back to deploy in parallel
	serviceIDs = make([]string, 0)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		volumes := d.extractVolumes(s)
		volumesFrom, err := d.extractVolumesFrom(s)
		if err != nil {
			return nil, err
		}
		serviceID, err := c.StartService(container.ServiceOptions{
			Namespace: d.namespace(s.namespace()),
			Labels: map[string]string{
				"mesg.service": s.Name,
				"mesg.hash":    s.Hash,
				"mesg.sid":     s.Sid,
				"mesg.core":    conf.Core.Name,
			},
			Image:   d.Image,
			Args:    d.Args,
			Command: d.Command,
			Env: xos.EnvMergeSlices(d.Env, []string{
				"MESG_TOKEN=" + s.Hash,
				"MESG_ENDPOINT=" + endpoint,
				"MESG_ENDPOINT_TCP=" + endpoint,
			}),
			Mounts: append(volumes, volumesFrom...),
			Ports:  d.extractPorts(),
			Networks: []container.Network{
				{ID: networkID, Alias: d.Key},
				{ID: sharedNetworkID},
			},
		})
		if err != nil {
			s.Stop(c)
			return nil, err
		}
		serviceIDs = append(serviceIDs, serviceID)
	}
	return serviceIDs, nil
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
func (d *Dependency) extractVolumes(s *Service) []container.Mount {
	volumes := make([]container.Mount, 0)
	for _, volume := range d.Volumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(s, d.Key, volume),
			Target: volume,
		})
	}
	return volumes
}

func (d *Dependency) extractVolumesFrom(s *Service) ([]container.Mount, error) {
	volumesFrom := make([]container.Mount, 0)
	for _, depName := range d.VolumesFrom {
		dep, err := s.getDependency(depName)
		if err != nil {
			if depName == MainServiceKey {
				dep = s.Configuration
			} else {
				return nil, err
			}
		}
		for _, volume := range dep.Volumes {
			volumesFrom = append(volumesFrom, container.Mount{
				Source: volumeKey(s, depName, volume),
				Target: volume,
			})
		}
	}
	return volumesFrom, nil
}

// volumeKey creates a key for service's volume based on the sid to make sure that the volume
// will stay the same for different versions of the service.
func volumeKey(s *Service, dependency string, volume string) string {
	return xstructhash.Hash([]string{
		s.Sid,
		dependency,
		volume,
	}, 1)
}
