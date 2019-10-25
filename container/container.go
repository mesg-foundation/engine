package container

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
)

// defaultStopGracePeriod is the timeout value between stopping a container and killing it.
const defaultStopGracePeriod = 10 * time.Second

var errSwarmNotInit = errors.New(`docker swarm is not initialized. run "docker swarm init" and try again`)

// Status of the service.
type Status uint

// Possible status for services.
const (
	UNKNOWN Status = iota
	STOPPED
	STARTING
	RUNNING
)

// statuses is a struct used to map service and contaienr statuses.
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

// buildResponse is the object that is returned by the docker build api.
type buildResponse struct {
	Stream string `json:"stream"`
	Error  string `json:"error"`
}

// Container describes the API of container package.
type Container interface {
	Build(path string) (tag string, err error)
	CreateNetwork(namespace string) (id string, err error)
	DeleteNetwork(namespace string) error
	SharedNetworkID() string
	StartService(options ServiceOptions) (serviceID string, err error)
	StopService(namespace string) (err error)
	DeleteVolume(name string) error
}

// DockerContainer provides high level interactions with Docker API for MESG.
type DockerContainer struct {
	// client is a Docker client.
	client client.CommonAPIClient

	// namespace prefix.
	nsprefix string

	// sharedNetworkID is cache for an id network.
	sharedNetworkID string
}

// New creates a new Container with given options.
func New(nsprefix string) (*DockerContainer, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	client.NegotiateAPIVersion(context.Background())

	c := &DockerContainer{
		client:   client,
		nsprefix: nsprefix,
	}

	if err := c.isSwarmInit(); err != nil {
		return nil, err
	}

	if err := c.createSharedNetwork(); err != nil {
		return nil, err
	}
	return c, nil
}

// Build builds a docker image.
func (c *DockerContainer) Build(path string) (tag string, err error) {
	excludeFiles, err := build.ReadDockerignore(path)
	if err != nil {
		return "", err
	}

	tr, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: excludeFiles,
	})
	if err != nil {
		return "", err
	}
	defer tr.Close()

	resp, err := c.client.ImageBuild(context.Background(), tr, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return tagFromResponse(resp.Body)
}

// CreateNetwork creates a Docker Network with a namespace. Retruns network id and error.
func (c *DockerContainer) CreateNetwork(namespace string) (string, error) {
	namespace = c.namespace(namespace)
	network, err := c.client.NetworkInspect(context.Background(), namespace, types.NetworkInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return "", err
	}
	if network.ID != "" {
		return network.ID, nil
	}
	response, err := c.client.NetworkCreate(context.Background(), namespace, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return response.ID, err
}

// DeleteNetwork deletes a Docker Network associated with a namespace.
func (c *DockerContainer) DeleteNetwork(namespace string) error {
	for {
		network, err := c.client.NetworkInspect(context.Background(), c.namespace(namespace), types.NetworkInspectOptions{})
		if client.IsErrNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
		c.client.NetworkRemove(context.Background(), network.ID)
		time.Sleep(100 * time.Millisecond)
	}
}

// SharedNetworkID returns the ID of the shared network.
func (c *DockerContainer) SharedNetworkID() string {
	return c.sharedNetworkID
}

// StartService starts a docker service.
func (c *DockerContainer) StartService(options ServiceOptions) (string, error) {
	if status, _ := c.serviceStatus(options.Namespace); status == RUNNING {
		service, _, err := c.client.ServiceInspectWithRaw(context.Background(), c.namespace(options.Namespace), types.ServiceInspectOptions{})
		return service.ID, err
	}

	service := options.toSwarmServiceSpec(c)
	response, err := c.client.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service.
func (c *DockerContainer) StopService(namespace string) error {
	status, err := c.serviceStatus(namespace)
	if err != nil {
		return err
	}
	if status == STOPPED {
		return nil
	}
	service, _, err := c.client.ServiceInspectWithRaw(context.Background(), c.namespace(namespace), types.ServiceInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	stopGracePeriod := defaultStopGracePeriod
	if service.Spec.TaskTemplate.ContainerSpec != nil &&
		service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod != nil {
		stopGracePeriod = *service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod
	}

	err = c.client.ServiceRemove(context.Background(), c.namespace(namespace))
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}
	if err := c.deletePendingContainer(namespace, time.Now().Add(stopGracePeriod)); err != nil {
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

// DeleteVolume deletes a Docker Volume by name.
func (c *DockerContainer) DeleteVolume(name string) error {
	return c.client.VolumeRemove(context.Background(), name, false)
}

// Cleanup cleans all configuration like shared network of running services.
func (c *DockerContainer) Cleanup() error {
	// remove shared network
	if err := c.client.NetworkRemove(context.Background(), c.nsprefix); err != nil {
		return err
	}
	return nil
}

// namespace creates a new namespace with container prefix.
func (c *DockerContainer) namespace(s string) string {
	return c.nsprefix + "-" + s
}

// isSwarmInit returns true if docker is connected with any swarm.
func (c *DockerContainer) isSwarmInit() error {
	info, err := c.client.Info(context.Background())
	if err != nil {
		return err
	}
	if info.Swarm.NodeID == "" {
		return errSwarmNotInit
	}
	return nil
}

// findContainer returns a docker container.
func (c *DockerContainer) findContainer(namespace string) (string, error) {
	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + c.namespace(namespace),
		}),
		Limit: 1,
	})
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", errdefs.NotFound(fmt.Errorf("container in namespace %s not found", namespace))
	}
	return containers[0].ID, nil
}

// Status returns the status of the container based on the docker container and docker service.
// if any error occurs during the status check, status will be shown as UNKNOWN.
// otherwise the following rules will be applied to determine a status:
//  - RUNNING: when the container is running in docker regardless of the status of the service.
//  - STARTING: when the service is running but the container is not yet started.
//  - STOPPED: when the container and the service is not running in docker.
func (c *DockerContainer) serviceStatus(namespace string) (Status, error) {
	container, err := c.containerExists(namespace)
	if err != nil {
		return UNKNOWN, err
	}

	service, err := c.serviceExists(namespace)
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

// containerExists checks if container with namespace can be found.
func (c *DockerContainer) containerExists(namespace string) (bool, error) {
	_, err := c.findContainer(namespace)
	if err != nil && !client.IsErrNotFound(err) {
		return false, err
	}
	return !client.IsErrNotFound(err), nil
}

// serviceExists checks if corresponding container for service namespace can be found.
func (c *DockerContainer) serviceExists(namespace string) (bool, error) {
	_, _, err := c.client.ServiceInspectWithRaw(context.Background(), c.namespace(namespace), types.ServiceInspectOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return false, err
	}
	return !client.IsErrNotFound(err), nil
}

// tasksError returns the error of matching tasks.
func (c *DockerContainer) tasksError(namespace string) ([]string, error) {
	tasks, err := c.client.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + c.namespace(namespace),
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
func (c *DockerContainer) waitForStatus(namespace string, status Status) error {
	curstatus, err := c.serviceStatus(namespace)
	if err != nil {
		return err
	}

	for curstatus != status {
		time.Sleep(100 * time.Millisecond)

		tasksErrors, err := c.tasksError(namespace)
		if err != nil {
			return err
		}
		if len(tasksErrors) > 0 {
			return errors.New(strings.Join(tasksErrors, ", "))
		}

		curstatus, err = c.serviceStatus(namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DockerContainer) createSharedNetwork() error {
	// check if already exist
	network, err := c.client.NetworkInspect(context.Background(), c.nsprefix, types.NetworkInspectOptions{})
	if network.ID != "" {
		c.sharedNetworkID = network.ID
		return nil
	}
	if err != nil && !client.IsErrNotFound(err) {
		return err
	}

	// Create the new network needed to run containers.
	resp, err := c.client.NetworkCreate(context.Background(), c.nsprefix, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": c.nsprefix,
		},
	})
	if err != nil {
		return err
	}

	c.sharedNetworkID = resp.ID
	return nil
}

func (c *DockerContainer) deletePendingContainer(namespace string, maxGraceTime time.Time) error {
	var (
		id  string
		err error
	)
	for start := time.Now(); start.Before(maxGraceTime); time.Sleep(100 * time.Millisecond) {
		id, err = c.findContainer(namespace)
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

// tagFromResponse retrives image build tag from client.Build response.
func tagFromResponse(r io.Reader) (string, error) {
	last, err := lastResponseLine(r)
	if err != nil {
		return "", err
	}

	var res buildResponse
	if err := json.Unmarshal(last, &res); err != nil {
		return "", fmt.Errorf("could not parse container build response: %s", err)
	}
	if res.Error != "" {
		return "", fmt.Errorf("image build failed: %s", res.Error)
	}
	if !strings.HasPrefix(res.Stream, "sha256:") {
		return "", fmt.Errorf("container: image build api dosen't return container id")
	}
	return strings.TrimSpace(res.Stream), nil
}

// lastResponseLine returns the last log line from client.Build response.
func lastResponseLine(r io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(bytes.TrimSpace(b), []byte{'\n'})
	if l := len(lines); l == 0 || len(lines[l-1]) == 0 {
		return nil, errors.New("container: image build api return empty response")
	}
	return lines[len(lines)-1], nil
}
