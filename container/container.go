package container

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xnet"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

const (
	// defaultStopGracePeriod is the timeout value between stopping a container and killing it.
	defaultStopGracePeriod = 10 * time.Second

	defaultMaxAttempts = uint64(3)

	servicePrefix = "mesg_srv_"
	imageTag      = "mesg:"

	pollingTime = 500 * time.Millisecond
)

// Status of the service.
type Status uint

// Possible status for services.
const (
	UNKNOWN Status = iota
	STOPPED
	STARTING
	RUNNING
)

// statuses is a struct used to map service and container statuses.
var statuses = []struct {
	container bool
	service   bool
	status    Status
}{
	{service: true, container: true, status: RUNNING},
	{service: true, container: false, status: STARTING},
	{service: false, container: true, status: RUNNING}, // This is actually stopping
	{service: false, container: false, status: STOPPED},
}

// Container starts and stops the MESG Service in Docker Container.
type Container struct {
	// client is a Docker client.
	client client.CommonAPIClient

	engineEndpoint string
	engineName     string
	engineNetwork  string
}

// New initializes the Container struct by connecting creating the Docker client.
func New(engineName, engineAddress, engineNetwork string) (*Container, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	client.NegotiateAPIVersion(context.Background())

	_, port, err := xnet.SplitHostPort(engineAddress)
	if err != nil {
		return nil, err
	}

	return &Container{
		client:         client,
		engineEndpoint: net.JoinHostPort(engineName, strconv.Itoa(port)),
		engineName:     engineName,
		engineNetwork:  engineNetwork,
	}, nil
}

// Download downloads the tarball of the source of a service from a HTTP url.
// Don't forget to remove the downloaded service.
func (c *Container) Download(url string) (string, error) {
	// download and untar service context into path.
	path, err := ioutil.TempDir("", "mesg")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("service's source code is not reachable, status: " + resp.Status + ", url: " + url)
	}
	defer resp.Body.Close()

	if err := archive.Untar(resp.Body, path, &archive.TarOptions{ChownOpts: &idtools.Identity{
		UID: os.Geteuid(),
		GID: os.Getegid()},
	}); err != nil {
		return "", err
	}

	return path, nil
}

// Build the imge of the container
func (c *Container) Build(srvHash hash.Hash, path string) error {
	excludeFiles, err := build.ReadDockerignore(path)
	if err != nil {
		return err
	}

	tr, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: excludeFiles,
	})
	if err != nil {
		return err
	}
	defer tr.Close()

	if _, err := c.client.ImageBuild(context.Background(), tr, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
		Tags:           []string{imageTag + srvHash.String()},
	}); err != nil {
		return err
	}

	return nil
}

// Start starts the service.
func (c *Container) Start(srv *service.Service, instanceHash, runnerHash, instanceEnvHash hash.Hash, instanceEnv []string, registerPayload []byte) (err error) {
	runnerName := servicePrefix + runnerHash.String()

	networkID, err := c.createNetwork(runnerName)
	if err != nil {
		return err
	}
	enginedNetworkID, err := c.enginedNetworkID()
	if err != nil {
		return err
	}

	specs := make([]swarm.ServiceSpec, 0)

	// Create dependency container configs
	for _, dep := range srv.Dependencies {
		depName := runnerName + "_" + dep.Key
		volumes := convertVolumes(srv, dep.Volumes, dep.Key)
		volumesFrom, err := convertVolumesFrom(srv, dep.VolumesFrom)
		if err != nil {
			return err
		}
		stopGracePeriod := defaultStopGracePeriod
		maxAttempts := defaultMaxAttempts
		specs = append(specs, swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: depName,
				Labels: map[string]string{
					"mesg.engine":                c.engineName,
					"mesg.service":               srv.Hash.String(),
					"mesg.instance":              instanceHash.String(),
					"mesg.runner":                runnerHash.String(),
					"mesg.dependency":            dep.Key,
					"com.docker.stack.namespace": depName,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: dep.Image,
					Labels: map[string]string{
						"com.docker.stack.namespace": depName,
					},
					Env:             dep.Env,
					Args:            dep.Args,
					Command:         strings.Fields(dep.Command),
					Mounts:          append(volumes, volumesFrom...),
					StopGracePeriod: &stopGracePeriod,
				},
				Networks: []swarm.NetworkAttachmentConfig{
					{
						Target:  networkID,
						Aliases: []string{dep.Key},
					},
				},
				RestartPolicy: &swarm.RestartPolicy{
					Condition:   swarm.RestartPolicyConditionOnFailure,
					MaxAttempts: &maxAttempts,
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: convertPorts(dep.Ports),
			},
		})
	}

	// Create configuration container config
	volumes := convertVolumes(srv, srv.Configuration.Volumes, service.MainServiceKey)
	volumesFrom, err := convertVolumesFrom(srv, srv.Configuration.VolumesFrom)
	if err != nil {
		return err
	}
	stopGracePeriod := defaultStopGracePeriod
	maxAttempts := defaultMaxAttempts
	specs = append(specs, swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: runnerName,
			Labels: map[string]string{
				"mesg.engine":                c.engineName,
				"mesg.service":               srv.Hash.String(),
				"mesg.instance":              instanceHash.String(),
				"mesg.runner":                runnerHash.String(),
				"mesg.dependency":            service.MainServiceKey,
				"com.docker.stack.namespace": runnerName,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: imageTag + srv.Hash.String(),
				Labels: map[string]string{
					"com.docker.stack.namespace": runnerName,
				},
				Args:    srv.Configuration.Args,
				Command: strings.Fields(srv.Configuration.Command),
				Env: xos.EnvMergeSlices(instanceEnv, []string{
					"MESG_SERVICE_HASH=" + srv.Hash.String(),
					"MESG_INSTANCE_HASH=" + instanceHash.String(),
					"MESG_RUNNER_HASH=" + runnerHash.String(),
					"MESG_ENDPOINT=" + c.engineEndpoint,
					"MESG_REGISTER_SIGNATURE=" + base64.StdEncoding.EncodeToString(registerPayload),
					"MESG_ENV_HASH=" + instanceEnvHash.String(),
				}),
				Mounts:          append(volumes, volumesFrom...),
				StopGracePeriod: &stopGracePeriod,
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  networkID,
					Aliases: []string{service.MainServiceKey},
				},
				{
					Target: enginedNetworkID,
				},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition:   swarm.RestartPolicyConditionOnFailure,
				MaxAttempts: &maxAttempts,
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: convertPorts(srv.Configuration.Ports),
		},
	})

	// Start
	for _, spec := range specs {
		if err := c.startService(spec); err != nil {
			c.Stop(srv, runnerHash)
			return err
		}
	}

	return nil
}

// Stop stops an instance.
func (c *Container) Stop(srv *service.Service, runnerHash hash.Hash) error {
	var (
		errs  xerrors.SyncErrors
		name  = servicePrefix + runnerHash.String()
		names = make([]string, 0)
	)

	for _, dep := range srv.Dependencies {
		names = append(names, name+"_"+dep.Key)
	}
	names = append(names, name)

	for _, name := range names {
		if err := c.stopService(name); err != nil {
			errs.Append(err)
		}
	}
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return c.deleteNetwork(name)
}

// deleteData deletes the data volumes of instance and its dependencies.
// TODO: right now deleteData() is not compatible to work with multiple instances of same
// service since extractVolumes() accepts service, not instance. same applies in the start
// api as well. make it work with multiple instances.
// func deleteData(cont container.Container, s *service.Service) error {
// 	var (
// 		wg      sync.WaitGroup
// 		errs    xerrors.SyncErrors
// 		volumes = make([]mount.Mount, 0)
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

func convertPorts(dPorts []string) []swarm.PortConfig {
	ports := make([]swarm.PortConfig, len(dPorts))
	for i, p := range dPorts {
		split := strings.Split(p, ":")
		from, _ := strconv.ParseUint(split[0], 10, 64)
		to := from
		if len(split) > 1 {
			to, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    uint32(to),
			PublishedPort: uint32(from),
		}
	}
	return ports
}

func convertVolumes(s *service.Service, dVolumes []string, key string) []mount.Mount {
	volumes := make([]mount.Mount, 0)
	for _, volume := range dVolumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeVolume,
			Source: volumeKey(s, key, volume),
			Target: volume,
		})
	}
	return volumes
}

func convertVolumesFrom(s *service.Service, dVolumesFrom []string) ([]mount.Mount, error) {
	volumesFrom := make([]mount.Mount, 0)
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
			volumesFrom = append(volumesFrom, mount.Mount{
				Type:   mount.TypeVolume,
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
		s.Hash.String(),
		dependency,
		volume,
	}).String()
}

// createNetwork creates a Docker Network with a name. Retruns network id and error.
func (c *Container) createNetwork(name string) (string, error) {
	network, err := c.client.NetworkInspect(context.Background(), name, types.NetworkInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return "", err
	}
	if network.ID != "" {
		return network.ID, nil
	}
	response, err := c.client.NetworkCreate(context.Background(), name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": name,
		},
	})
	return response.ID, err
}

// deleteNetwork deletes a Docker Network.
func (c *Container) deleteNetwork(name string) error {
	for {
		network, err := c.client.NetworkInspect(context.Background(), name, types.NetworkInspectOptions{})
		if client.IsErrNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
		c.client.NetworkRemove(context.Background(), network.ID)
		time.Sleep(pollingTime)
	}
}

// enginedNetworkID retrieve the docker network id of the engine network.
func (c *Container) enginedNetworkID() (string, error) {
	network, err := c.client.NetworkInspect(context.Background(), c.engineNetwork, types.NetworkInspectOptions{})
	return network.ID, err
}

// StartService starts a docker service.
func (c *Container) startService(spec swarm.ServiceSpec) error {
	if status, _ := c.serviceStatus(spec.Name); status == RUNNING {
		_, _, err := c.client.ServiceInspectWithRaw(context.Background(), spec.Name, types.ServiceInspectOptions{})
		return err
	}
	if _, err := c.client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{}); err != nil {
		return err
	}
	return c.waitForStatus(spec.Name, RUNNING)
}

// StopService stops a docker service.
func (c *Container) stopService(name string) error {
	status, err := c.serviceStatus(name)
	if err != nil {
		return err
	}
	if status == STOPPED {
		return nil
	}
	service, _, err := c.client.ServiceInspectWithRaw(context.Background(), name, types.ServiceInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	stopGracePeriod := defaultStopGracePeriod
	if service.Spec.TaskTemplate.ContainerSpec != nil &&
		service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod != nil {
		stopGracePeriod = *service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod
	}

	err = c.client.ServiceRemove(context.Background(), name)
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	if err := c.deletePendingContainer(name, time.Now().Add(stopGracePeriod)); err != nil {
		return err
	}
	return c.waitForStatus(name, STOPPED)
}

// findContainer returns a docker container.
func (c *Container) findContainer(name string) (string, error) {
	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + name,
		}),
		Limit: 1,
	})
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", errdefs.NotFound(fmt.Errorf("container in namespace %s not found", name))
	}
	return containers[0].ID, nil
}

// Status returns the status of the container based on the docker container and docker service.
// if any error occurs during the status check, status will be shown as UNKNOWN.
// otherwise the following rules will be applied to determine a status:
//  - RUNNING: when the container is running in docker regardless of the status of the service.
//  - STARTING: when the service is running but the container is not yet started.
//  - STOPPED: when the container and the service is not running in docker.
func (c *Container) serviceStatus(name string) (Status, error) {
	container, err := c.containerExists(name)
	if err != nil {
		return UNKNOWN, err
	}

	service, err := c.serviceExists(name)
	if err != nil {
		return UNKNOWN, err
	}

	for _, s := range statuses {
		if s.container == container && s.service == service {
			return s.status, nil
		}
	}

	// This should never be reached but it's better than a panic :)
	return UNKNOWN, nil
}

// containerExists checks if container with name can be found.
func (c *Container) containerExists(name string) (bool, error) {
	_, err := c.findContainer(name)
	if err != nil && !client.IsErrNotFound(err) {
		return false, err
	}
	return !client.IsErrNotFound(err), nil
}

// serviceExists checks if corresponding container for service namespace can be found.
func (c *Container) serviceExists(name string) (bool, error) {
	_, _, err := c.client.ServiceInspectWithRaw(context.Background(), name, types.ServiceInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return false, err
	}
	return !client.IsErrNotFound(err), nil
}

// tasksError returns the error of matching tasks.
func (c *Container) tasksError(name string) ([]string, error) {
	tasks, err := c.client.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + name,
		}),
	})
	if err != nil {
		return nil, err
	}

	var errors []string
	for _, task := range tasks {
		if task.Status.Err != "" {
			errors = append(errors, task.Status.Err)
		}
	}
	return errors, nil
}

// waitForStatus waits for the container to have the provided status. Returns error as soon as possible.
func (c *Container) waitForStatus(name string, status Status) error {
	curstatus, err := c.serviceStatus(name)
	if err != nil {
		return err
	}

	for curstatus != status {
		time.Sleep(pollingTime)

		tasksErrors, err := c.tasksError(name)
		if err != nil {
			return err
		}
		if len(tasksErrors) > 0 {
			return errors.New(strings.Join(tasksErrors, ", "))
		}

		curstatus, err = c.serviceStatus(name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) deletePendingContainer(name string, maxGraceTime time.Time) error {
	var (
		id  string
		err error
	)
	for start := time.Now(); start.Before(maxGraceTime); time.Sleep(pollingTime) {
		id, err = c.findContainer(name)
		if client.IsErrNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
	}
	// Hack to force Docker to remove the containers.
	// Sometime, the ServiceRemove function doesn't remove the associated containers,
	// or too late and the associated networks cannot be removed.
	// This hack for Docker to stop and then remove the container.
	// See issue https://github.com/moby/moby/issues/32620
	c.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true})
	return nil
}
