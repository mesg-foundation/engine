package instancesdk

import (
	"strconv"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/mesg-foundation/core/x/xos"
)

// Start starts the service.
func (i *Instance) start(inst *instance.Instance) (serviceIDs []string, err error) {
	srv, err := i.serviceDB.Get(inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	sNamespace := instanceNamespace(inst.Hash)
	networkID, err := i.container.CreateNetwork(sNamespace)
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
	endpoint := conf.Name + ":" + strconv.Itoa(port)
	env := xos.EnvMapToSlice(xos.EnvMergeMaps(xos.EnvSliceToMap(srv.Configuration.Env), xos.EnvSliceToMap(inst.Env)))
	// BUG: https://github.com/mesg-foundation/core/issues/382
	// After solving this by docker, switch back to deploy in parallel
	serviceIDs = make([]string, 0)
	for _, d := range append(srv.Dependencies, srv.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		volumes := extractVolumes(srv, d)
		volumesFrom, err := extractVolumesFrom(srv, d)
		if err != nil {
			return nil, err
		}
		serviceID, err := i.container.StartService(container.ServiceOptions{
			Namespace: dependencyNamespace(sNamespace, d.Key),
			Labels: map[string]string{
				"mesg.service": srv.Name,
				"mesg.hash":    inst.Hash,
				"mesg.sid":     srv.Sid,
				"mesg.engine":  conf.Name,
			},
			Image:   d.Image,
			Args:    d.Args,
			Command: d.Command,
			Env: xos.EnvMergeSlices(env, []string{
				"MESG_TOKEN=" + inst.Hash,
				"MESG_ENDPOINT=" + endpoint,
				"MESG_ENDPOINT_TCP=" + endpoint,
			}),
			Mounts: append(volumes, volumesFrom...),
			Ports:  extractPorts(d),
			Networks: []container.Network{
				{ID: networkID, Alias: d.Key},
				{ID: sharedNetworkID},
			},
		})
		if err != nil {

			i.stop(inst)
			return nil, err
		}
		serviceIDs = append(serviceIDs, serviceID)
	}
	return serviceIDs, nil
}
