package instancesdk

import (
	"strconv"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xos"
)

// Start starts the service.
func (i *Instance) start(inst *instance.Instance, imageHash string, env []string) (serviceIDs []string, err error) {
	srv, err := i.service.Get(inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	instNamespace := InstanceNamespace(inst.Hash)
	networkID, err := i.container.CreateNetwork(instNamespace)
	if err != nil {
		return nil, err
	}
	sharedNetworkID, err := i.container.SharedNetworkID()
	if err != nil {
		return nil, err
	}
	conf, err := config.Global()
	if err != nil {
		return nil, err
	}
	_, port, _ := xnet.SplitHostPort(conf.Server.Address)
	endpoint := version.Name + ":" + strconv.Itoa(port)
	// BUG: https://github.com/mesg-foundation/engine/issues/382
	// After solving this by docker, switch back to deploy in parallel
	configs := make([]container.ServiceOptions, 0)

	labels := map[string]string{
		"mesg.service": srv.Name,
		"mesg.hash":    inst.Hash.String(),
		"mesg.sid":     srv.Sid,
		"mesg.engine":  version.Name,
	}

	// Create dependency container configs
	for _, d := range srv.Dependencies {
		volumes := extractVolumes(srv, d)
		volumesFrom, err := extractVolumesFrom(srv, d)
		if err != nil {
			return nil, err
		}
		configs = append(configs, container.ServiceOptions{
			Namespace: DependencyNamespace(instNamespace, d.Key),
			Labels:    labels,
			Image:     d.Image,
			Args:      d.Args,
			Command:   d.Command,
			Env:       d.Env,
			Mounts:    append(volumes, volumesFrom...),
			Ports:     extractPorts(d),
			Networks: []container.Network{
				{ID: networkID, Alias: d.Key},
				{ID: sharedNetworkID}, // TODO: to remove
			},
		})
	}

	// Create configuration container config
	volumes := extractVolumes(srv, srv.Configuration)
	volumesFrom, err := extractVolumesFrom(srv, srv.Configuration)
	if err != nil {
		return nil, err
	}
	configs = append(configs, container.ServiceOptions{
		Namespace: DependencyNamespace(instNamespace, srv.Configuration.Key),
		Labels:    labels,
		Image:     imageHash,
		Args:      srv.Configuration.Args,
		Command:   srv.Configuration.Command,
		Env: xos.EnvMergeSlices(env, []string{
			"MESG_TOKEN=" + inst.Hash.String(),
			"MESG_ENDPOINT=" + endpoint,
			"MESG_ENDPOINT_TCP=" + endpoint,
		}),
		Mounts: append(volumes, volumesFrom...),
		Ports:  extractPorts(srv.Configuration),
		Networks: []container.Network{
			{ID: networkID, Alias: srv.Configuration.Key},
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
