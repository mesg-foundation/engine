package builder

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

// Build the imge of the container
func Build(cont container.Container, srv *service.Service, ipfsEndpoint string) (string, error) {
	// download and untar service context into path.
	path, err := ioutil.TempDir("", "mesg")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(path)

	resp, err := http.Get(ipfsEndpoint + srv.Source)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("service's source code is not reachable, status: " + resp.Status + ", url: " + ipfsEndpoint + srv.Source)
	}
	defer resp.Body.Close()

	if err := archive.Untar(resp.Body, path, &archive.TarOptions{ChownOpts: &idtools.Identity{
		UID: os.Geteuid(),
		GID: os.Getegid()},
	}); err != nil {
		return "", err
	}

	// build service's Docker image and apply to service.
	imageHash, err := cont.Build(path)
	if err != nil {
		return "", err
	}

	return imageHash, nil
}

// Start starts the service.
func Start(cont container.Container, srv *service.Service, instanceHash hash.Hash, runnerHash hash.Hash, imageHash string, env []string, engineName, port string) (err error) {
	endpoint := net.JoinHostPort(engineName, port)
	namespace := namespace(runnerHash)
	networkID, err := cont.CreateNetwork(namespace)
	if err != nil {
		return err
	}
	sharedNetworkID := cont.SharedNetworkID()
	// BUG: https://github.com/mesg-foundation/engine/issues/382
	// After solving this by docker, switch back to deploy in parallel
	configs := make([]container.ServiceOptions, 0)

	// Create dependency container configs
	for _, d := range srv.Dependencies {
		volumes := convertVolumes(srv, d.Volumes, d.Key)
		volumesFrom, err := convertVolumesFrom(srv, d.VolumesFrom)
		if err != nil {
			return err
		}
		configs = append(configs, container.ServiceOptions{
			Namespace: dependencyNamespace(namespace, d.Key),
			Labels: map[string]string{
				"mesg.engine":     engineName,
				"mesg.sid":        srv.Sid,
				"mesg.service":    srv.Hash.String(),
				"mesg.instance":   instanceHash.String(),
				"mesg.runner":     runnerHash.String(),
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
		return err
	}
	configs = append(configs, container.ServiceOptions{
		Namespace: dependencyNamespace(namespace, service.MainServiceKey),
		Labels: map[string]string{
			"mesg.engine":     engineName,
			"mesg.sid":        srv.Sid,
			"mesg.service":    srv.Hash.String(),
			"mesg.instance":   instanceHash.String(),
			"mesg.runner":     runnerHash.String(),
			"mesg.dependency": service.MainServiceKey,
		},
		Image:   imageHash,
		Args:    srv.Configuration.Args,
		Command: srv.Configuration.Command,
		Env: xos.EnvMergeSlices(env, []string{
			"MESG_SERVICE_HASH=" + srv.Hash.String(),
			"MESG_INSTANCE_HASH=" + instanceHash.String(),
			"MESG_RUNNER_HASH=" + runnerHash.String(),
			"MESG_ENDPOINT=" + endpoint,
		}),
		Mounts: append(volumes, volumesFrom...),
		Ports:  convertPorts(srv.Configuration.Ports),
		Networks: []container.Network{
			{ID: networkID, Alias: service.MainServiceKey},
			{ID: sharedNetworkID},
		},
	})

	// Start
	for _, c := range configs {
		_, err := cont.StartService(c)
		if err != nil {
			Stop(cont, runnerHash, srv.Dependencies)
			return err
		}
	}

	return nil
}

// Stop stops an instance.
func Stop(cont container.Container, runnerHash hash.Hash, dependencies []*service.Service_Dependency) error {
	var (
		wg         sync.WaitGroup
		errs       xerrors.SyncErrors
		namespace  = namespace(runnerHash)
		namespaces = make([]string, 0)
	)

	for _, d := range dependencies {
		namespaces = append(namespaces, dependencyNamespace(namespace, d.Key))
	}
	namespaces = append(namespaces, dependencyNamespace(namespace, service.MainServiceKey))

	for _, namespace := range namespaces {
		wg.Add(1)
		go func(namespace string) {
			defer wg.Done()
			if err := cont.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}(namespace)
	}
	wg.Wait()
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return cont.DeleteNetwork(namespace)
}

// deleteData deletes the data volumes of instance and its dependencies.
// TODO: right now deleteData() is not compatible to work with multiple instances of same
// service since extractVolumes() accepts service, not instance. same applies in the start
// api as well. make it work with multiple instances.
// func deleteData(cont container.Container, s *service.Service) error {
// 	var (
// 		wg      sync.WaitGroup
// 		errs    xerrors.SyncErrors
// 		volumes = make([]container.Mount, 0)
// 	)

// 	for _, d := range s.Dependencies {
// 		volumes = append(volumes, convertVolumes(s, d.Volumes, d.Key)...)
// 	}
// 	volumes = append(volumes, convertVolumes(s, s.Configuration.Volumes, service.MainServiceKey)...)

// 	for _, volume := range volumes {
// 		// TODO: is it actually necessary to remvoe in parallel? I worry that some volume will be deleted at the same time and create side effect
// 		wg.Add(1)
// 		go func(source string) {
// 			defer wg.Done()
// 			// if service is never started before, data volume won't be created and Docker Engine
// 			// will return with the not found error. therefore, we can safely ignore it.
// 			if err := cont.DeleteVolume(source); err != nil && !client.IsErrNotFound(err) {
// 				errs.Append(err)
// 			}
// 		}(volume.Source)
// 	}
// 	wg.Wait()
// 	return errs.ErrorOrNil()
// }

// namespace returns the namespace of the service.
func namespace(hash hash.Hash) string {
	return hash.String()
}

// dependencyNamespace builds the namespace of a dependency.
func dependencyNamespace(namespace string, dependencyKey string) string {
	return hash.Dump(namespace + dependencyKey).String()
}

func convertPorts(dPorts []string) []container.Port {
	ports := make([]container.Port, len(dPorts))
	for i, p := range dPorts {
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
func convertVolumes(s *service.Service, dVolumes []string, key string) []container.Mount {
	volumes := make([]container.Mount, 0)
	for _, volume := range dVolumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(s, key, volume),
			Target: volume,
		})
	}
	return volumes
}

func convertVolumesFrom(s *service.Service, dVolumesFrom []string) ([]container.Mount, error) {
	volumesFrom := make([]container.Mount, 0)
	for _, depName := range dVolumesFrom {
		var depVolumes []string
		if depName == service.MainServiceKey {
			depVolumes = s.Configuration.Volumes
		} else {
			dep, err := s.GetDependency(depName)
			if err != nil {
				return nil, err
			}
			depVolumes = dep.Volumes
		}
		for _, volume := range depVolumes {
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
func volumeKey(s *service.Service, dependency, volume string) string {
	return hash.Dump([]string{
		s.Sid,
		dependency,
		volume,
	}).String()
}
