package instancesdk

import (
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/xos"
)

// Start starts the service.
func (i *Instance) start(inst *instance.Instance, imageHash string, env []string) (serviceIDs []string, err error) {
	srv, err := i.service.Get(inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	instNamespace := instanceNamespace(inst.Hash)
	networkID, err := i.container.CreateNetwork(instNamespace)
	if err != nil {
		return nil, err
	}
	sharedNetworkID := i.container.SharedNetworkID()
	// BUG: https://github.com/mesg-foundation/engine/issues/382
	// After solving this by docker, switch back to deploy in parallel
	configs := make([]container.ServiceOptions, 0)

	// Create dependency container configs
	for _, d := range srv.Dependencies {
		volumes := convertVolumes(srv, d.Volumes, d.Key)
		volumesFrom, err := convertVolumesFrom(srv, d.VolumesFrom)
		if err != nil {
			return nil, err
		}
		configs = append(configs, container.ServiceOptions{
			Namespace: dependencyNamespace(instNamespace, d.Key),
			Labels: map[string]string{
				"mesg.engine":     i.engineName,
				"mesg.sid":        srv.Sid,
				"mesg.service":    srv.Hash.String(),
				"mesg.instance":   inst.Hash.String(),
				"mesg.dependency": d.Key,
			},
			Image:   d.Image,
			Args:    d.Args,
			Command: d.Command,
			Env:     d.Env,
			Mounts:  append(volumes, volumesFrom...),
			Ports:   convertPorts(d.Ports),
			Networks: []container.Network{
				{ID: networkID, Alias: d.Key},
				{ID: sharedNetworkID}, // TODO: to remove
			},
		})
	}

	// Create configuration container config
	volumes := convertVolumes(srv, srv.Configuration.Volumes, service.MainServiceKey)
	volumesFrom, err := convertVolumesFrom(srv, srv.Configuration.VolumesFrom)
	if err != nil {
		return nil, err
	}
	configs = append(configs, container.ServiceOptions{
		Namespace: dependencyNamespace(instNamespace, service.MainServiceKey),
		Labels: map[string]string{
			"mesg.engine":     i.engineName,
			"mesg.sid":        srv.Sid,
			"mesg.service":    srv.Hash.String(),
			"mesg.instance":   inst.Hash.String(),
			"mesg.dependency": service.MainServiceKey,
		},
		Image:   imageHash,
		Args:    srv.Configuration.Args,
		Command: srv.Configuration.Command,
		Env: xos.EnvMergeSlices(env, []string{
			"MESG_TOKEN=" + inst.Hash.String(),
			"MESG_INSTANCE_HASH=" + inst.Hash.String(),
			"MESG_ENDPOINT=" + i.endpoint,
		}),
		Mounts: append(volumes, volumesFrom...),
		Ports:  convertPorts(srv.Configuration.Ports),
		Networks: []container.Network{
			{ID: networkID, Alias: service.MainServiceKey},
			{ID: sharedNetworkID},
		},
	})

	// Start
	serviceIDs = make([]string, 0)
	for _, c := range configs {
		serviceID, err := i.container.StartService(c)
		if err != nil {
			i.stop(inst)
			return nil, err
		}
		serviceIDs = append(serviceIDs, serviceID)
	}

	return serviceIDs, nil
}
