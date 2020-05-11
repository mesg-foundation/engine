package container

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/go-connections/nat"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xnet"
	"github.com/mesg-foundation/engine/ext/xos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

const (
	// DefaultStopGracePeriod is the default timeout value between stopping a container and killing it.
	DefaultStopGracePeriod = 10 * time.Second

	// DefaultMaximumRetryCount is the default number of time to restart the container if it crashes.
	DefaultMaximumRetryCount = 3

	servicePrefix = "mesg_srv_"
	imageTag      = "mesg:"
)

// Container starts and stops the MESG Service in Docker Container.
type Container struct {
	client            client.CommonAPIClient
	engineEndpoint    string
	engineName        string
	engineNetwork     string
	stopGracePeriod   time.Duration
	maximumRetryCount int
}

// New initializes the Container struct by connecting creating the Docker client.
func New(engineName, engineAddress, engineNetwork string, maximumRetryCount int, stopGracePeriod time.Duration) (*Container, error) {
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
		client:            client,
		engineEndpoint:    net.JoinHostPort(engineName, strconv.Itoa(port)),
		engineName:        engineName,
		engineNetwork:     engineNetwork,
		stopGracePeriod:   stopGracePeriod,
		maximumRetryCount: maximumRetryCount,
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
	// delete the service's container on any error
	errorOccurred := true
	defer func() {
		if errorOccurred {
			c.Stop(srv, runnerHash)
		}
	}()

	runnerName := servicePrefix + runnerHash.String()

	networkID, err := c.createNetwork(runnerName)
	if err != nil {
		return err
	}
	enginedNetworkID, err := c.enginedNetworkID()
	if err != nil {
		return err
	}

	// Create dependency container configs
	for _, dep := range srv.Dependencies {
		depName := runnerName + "_" + dep.Key
		volumes := convertVolumes(srv, dep.Volumes, dep.Key)
		volumesFrom, err := convertVolumesFrom(srv, dep.VolumesFrom)
		if err != nil {
			return err
		}
		exposedPort, portBindings, err := nat.ParsePortSpecs(dep.Ports)
		containerConfig := &containertypes.Config{
			Image: dep.Image,
			Labels: map[string]string{
				"mesg.engine":     c.engineName,
				"mesg.service":    srv.Hash.String(),
				"mesg.instance":   instanceHash.String(),
				"mesg.runner":     runnerHash.String(),
				"mesg.dependency": dep.Key,
			},
			Env:          dep.Env,
			ExposedPorts: exposedPort,
		}
		hostConfig := &containertypes.HostConfig{
			PortBindings: portBindings,
			Mounts:       append(volumes, volumesFrom...),
			RestartPolicy: containertypes.RestartPolicy{
				Name:              "on-failure",
				MaximumRetryCount: c.maximumRetryCount,
			},
		}
		if dep.Command != "" {
			containerConfig.Cmd = append(containerConfig.Cmd, dep.Command)
		}
		containerConfig.Cmd = append(containerConfig.Cmd, dep.Args...)

		// create container, attach network, start it
		if _, _, err := c.client.ImageInspectWithRaw(context.Background(), dep.Image); err != nil {
			r, err := c.client.ImagePull(context.Background(), dep.Image, types.ImagePullOptions{})
			if err != nil {
				return err
			}
			// wait for docker to download the image
			if _, err := ioutil.ReadAll(r); err != nil {
				return err
			}
		}
		if _, err := c.client.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, depName); err != nil {
			return err
		}
		if err := c.client.NetworkConnect(context.Background(), networkID, depName, &networktypes.EndpointSettings{
			Aliases: []string{dep.Key},
		}); err != nil {
			return err
		}
		if err := c.client.ContainerStart(context.Background(), depName, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}

	// Create configuration container config
	exposedPort, portBindings, err := nat.ParsePortSpecs(srv.Configuration.Ports)
	if err != nil {
		return err
	}
	volumes := convertVolumes(srv, srv.Configuration.Volumes, service.MainServiceKey)
	volumesFrom, err := convertVolumesFrom(srv, srv.Configuration.VolumesFrom)
	if err != nil {
		return err
	}
	containerConfig := &containertypes.Config{
		Image: imageTag + srv.Hash.String(),
		Labels: map[string]string{
			"mesg.engine":     c.engineName,
			"mesg.service":    srv.Hash.String(),
			"mesg.instance":   instanceHash.String(),
			"mesg.runner":     runnerHash.String(),
			"mesg.dependency": service.MainServiceKey,
		},
		Env: xos.EnvMergeSlices(instanceEnv, []string{
			"MESG_SERVICE_HASH=" + srv.Hash.String(),
			"MESG_INSTANCE_HASH=" + instanceHash.String(),
			"MESG_RUNNER_HASH=" + runnerHash.String(),
			"MESG_ENDPOINT=" + c.engineEndpoint,
			"MESG_REGISTER_SIGNATURE=" + base64.StdEncoding.EncodeToString(registerPayload),
			"MESG_ENV_HASH=" + instanceEnvHash.String(),
		}),
		ExposedPorts: exposedPort,
	}
	hostConfig := &containertypes.HostConfig{
		PortBindings: portBindings,
		Mounts:       append(volumes, volumesFrom...),
		RestartPolicy: containertypes.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: c.maximumRetryCount,
		},
	}
	if srv.Configuration.Command != "" {
		containerConfig.Cmd = append(containerConfig.Cmd, srv.Configuration.Command)
	}
	containerConfig.Cmd = append(containerConfig.Cmd, srv.Configuration.Args...)

	// create container, attach networks, start it
	if _, err := c.client.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, runnerName); err != nil {
		return err
	}
	if err := c.client.NetworkConnect(context.Background(), networkID, runnerName, &networktypes.EndpointSettings{
		Aliases: []string{service.MainServiceKey},
	}); err != nil {
		return err
	}
	if err := c.client.NetworkConnect(context.Background(), enginedNetworkID, runnerName, &networktypes.EndpointSettings{}); err != nil {
		return err
	}
	if err := c.client.ContainerStart(context.Background(), runnerName, types.ContainerStartOptions{}); err != nil {
		return err
	}

	errorOccurred = false
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
		stopGracePeriod := c.stopGracePeriod
		if err := c.client.ContainerStop(context.Background(), name, &stopGracePeriod); err != nil {
			errs.Append(err)
		}
		if err := c.client.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{}); err != nil {
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

func convertVolumes(s *service.Service, dVolumes []string, key string) []mounttypes.Mount {
	volumes := make([]mounttypes.Mount, 0)
	for _, volume := range dVolumes {
		volumes = append(volumes, mounttypes.Mount{
			Type:   mounttypes.TypeVolume,
			Source: volumeKey(s, key, volume),
			Target: volume,
		})
	}
	return volumes
}

func convertVolumesFrom(s *service.Service, dVolumesFrom []string) ([]mounttypes.Mount, error) {
	volumesFrom := make([]mounttypes.Mount, 0)
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
			volumesFrom = append(volumesFrom, mounttypes.Mount{
				Type:   mounttypes.TypeVolume,
				Source: volumeKey(s, depName, volume),
				Target: volume,
			})
		}
	}
	return volumesFrom, nil
}

func volumeKey(s *service.Service, dependency, volume string) string {
	return hash.Dump([]string{
		s.Hash.String(),
		dependency,
		volume,
	}).String()
}

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
		Driver:         "bridge",
	})
	return response.ID, err
}

func (c *Container) deleteNetwork(name string) error {
	network, err := c.client.NetworkInspect(context.Background(), name, types.NetworkInspectOptions{})
	if client.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return c.client.NetworkRemove(context.Background(), network.ID)
}

func (c *Container) enginedNetworkID() (string, error) {
	network, err := c.client.NetworkInspect(context.Background(), c.engineNetwork, types.NetworkInspectOptions{})
	return network.ID, err
}
