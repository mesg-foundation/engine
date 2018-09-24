package dockermock

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/mock"
)

// CommonAPIClientMock is the common methods between stable and experimental versions of APIClient.
type CommonAPIClientMock struct {
	*mock.Mock
	*ConfigAPIClientMock
	*ContainerAPIClientMock
	*DistributionAPIClientMock
	*ImageAPIClientMock
	*NodeAPIClientMock
	*NetworkAPIClientMock
	*PluginAPIClientMock
	*ServiceAPIClientMock
	*SwarmAPIClientMock
	*SecretAPIClientMock
	*SystemAPIClientMock
	*VolumeAPIClientMock
}

// NewCommonAPIClientMock creates docker ContainerAPIClient mock.
func NewCommonAPIClientMock() *CommonAPIClientMock {
	m := &mock.Mock{}
	return &CommonAPIClientMock{
		Mock:                      m,
		ConfigAPIClientMock:       &ConfigAPIClientMock{m},
		ContainerAPIClientMock:    &ContainerAPIClientMock{m},
		DistributionAPIClientMock: &DistributionAPIClientMock{m},
		ImageAPIClientMock:        &ImageAPIClientMock{m},
		NodeAPIClientMock:         &NodeAPIClientMock{m},
		NetworkAPIClientMock:      &NetworkAPIClientMock{m},
		PluginAPIClientMock:       &PluginAPIClientMock{m},
		ServiceAPIClientMock:      &ServiceAPIClientMock{m},
		SwarmAPIClientMock:        &SwarmAPIClientMock{m},
		SecretAPIClientMock:       &SecretAPIClientMock{m},
		SystemAPIClientMock:       &SystemAPIClientMock{m},
		VolumeAPIClientMock:       &VolumeAPIClientMock{m},
	}
}

// ClientVersion returns the API version used by this client.
func (m *CommonAPIClientMock) ClientVersion() string {
	args := m.Called()
	return args.String(0)
}

// DaemonHost returns the host address used by the client.
func (m *CommonAPIClientMock) DaemonHost() string {
	args := m.Called()
	return args.String(0)
}

// HTTPClient returns a copy of the HTTP client bound to the server.
func (m *CommonAPIClientMock) HTTPClient() *http.Client {
	args := m.Called()
	return args.Get(0).(*http.Client)
}

// ServerVersion returns information of the docker client and server host.
func (m *CommonAPIClientMock) ServerVersion(ctx context.Context) (types.Version, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Version), args.Error(1)
}

// NegotiateAPIVersion queries the API and updates the version to match the API
// version. Any errors are silently ignored.
func (m *CommonAPIClientMock) NegotiateAPIVersion(ctx context.Context) {
	m.Called(ctx)
}

// NegotiateAPIVersionPing updates the client version to match the
// Ping.APIVersion if the ping version is less than the default version.
func (m *CommonAPIClientMock) NegotiateAPIVersionPing(ping types.Ping) {
	m.Called(ping)
}

// DialSession returns a connection that can be used communication with daemon.
func (m *CommonAPIClientMock) DialSession(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
	args := m.Called(ctx, proto, meta)
	return args.Get(0).(net.Conn), args.Error(1)
}

// Dialer returns a dialer for a raw stream connection, with HTTP/1.1 header,
// that can be used for proxying the daemon connection. Used by `docker
// dial-stdio`.
func (m *CommonAPIClientMock) Dialer() func(context.Context) (net.Conn, error) {
	args := m.Called()
	return args.Get(0).(func(context.Context) (net.Conn, error))
}

// Close the transport used by the client.
func (m *CommonAPIClientMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ContainerAPIClientMock defines API client methods for the containers
type ContainerAPIClientMock struct {
	*mock.Mock
}

// ContainerAttach attaches a connection to a container in the server. It
// returns a types.HijackedConnection with the hijacked connection and the a
// reader to get output. It's up to the called to close the hijacked connection
// by calling types.HijackedResponse.Close.
func (m *ContainerAPIClientMock) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(types.HijackedResponse), args.Error(1)
}

// ContainerCommit applies changes into a container and creates a new tagged
// image.
func (m *ContainerAPIClientMock) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (types.IDResponse, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(types.IDResponse), args.Error(1)
}

// ContainerCreate creates a new container based in the given configuration. It
// can be associated with a name, but it's not mandatory.
func (m *ContainerAPIClientMock) ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, networkingConfig *networktypes.NetworkingConfig, containerName string) (containertypes.ContainerCreateCreatedBody, error) {
	args := m.Called(ctx, config, hostConfig, networkingConfig, containerName)
	return args.Get(0).(containertypes.ContainerCreateCreatedBody), args.Error(1)
}

// ContainerDiff shows differences in a container filesystem since it was
// started.
func (m *ContainerAPIClientMock) ContainerDiff(ctx context.Context, container string) ([]containertypes.ContainerChangeResponseItem, error) {
	args := m.Called(ctx, container)
	return args.Get(0).([]containertypes.ContainerChangeResponseItem), args.Error(1)
}

// ContainerExecAttach attaches a connection to an exec process in the server.
// It returns a types.HijackedConnection with the hijacked connection and the a
// reader to get output. It's up to the called to close the hijacked connection
// by calling types.HijackedResponse.Close.
func (m *ContainerAPIClientMock) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	args := m.Called(ctx, execID, config)
	return args.Get(0).(types.HijackedResponse), args.Error(1)
}

// ContainerExecCreate creates a new exec configuration to run an exec process.
func (m *ContainerAPIClientMock) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	args := m.Called(ctx, container, config)
	return args.Get(0).(types.IDResponse), args.Error(1)
}

// ContainerExecInspect returns information about a specific exec process on
// the docker host.
func (m *ContainerAPIClientMock) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	args := m.Called(ctx, execID)
	return args.Get(0).(types.ContainerExecInspect), args.Error(1)
}

// ContainerExecResize changes the size of the tty for an exec process running
// inside a container.
func (m *ContainerAPIClientMock) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	args := m.Called(ctx, execID, options)
	return args.Error(0)
}

// ContainerExecStart starts an exec process already created in the docker
// host.
func (m *ContainerAPIClientMock) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	args := m.Called(ctx, execID, config)
	return args.Error(0)
}

// ContainerExport retrieves the raw contents of a container and returns them
// as an io.ReadCloser. It's up to the caller to close the stream.
func (m *ContainerAPIClientMock) ContainerExport(ctx context.Context, container string) (io.ReadCloser, error) {
	args := m.Called(ctx, container)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ContainerInspect returns the container information.
func (m *ContainerAPIClientMock) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	args := m.Called(ctx, container)
	return args.Get(0).(types.ContainerJSON), args.Error(1)
}

// ContainerInspectWithRaw returns the container information and its raw
// representation.
func (m *ContainerAPIClientMock) ContainerInspectWithRaw(ctx context.Context, container string, getSize bool) (types.ContainerJSON, []byte, error) {
	args := m.Called(ctx, container, getSize)
	return args.Get(0).(types.ContainerJSON), args.Get(1).([]byte), args.Error(2)
}

// ContainerKill terminates the container process but does not remove the
// container from the docker host.
func (m *ContainerAPIClientMock) ContainerKill(ctx context.Context, container, signal string) error {
	args := m.Called(ctx, container, signal)
	return args.Error(0)
}

// ContainerList returns the list of containers in the docker host.
func (m *ContainerAPIClientMock) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}

// ContainerLogs returns the logs generated by a container in an io.ReadCloser.
// It's up to the caller to close the stream.
func (m *ContainerAPIClientMock) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ContainerPause pauses the main process of a given container without
// terminating it.
func (m *ContainerAPIClientMock) ContainerPause(ctx context.Context, container string) error {
	args := m.Called(ctx, container)
	return args.Error(0)
}

// ContainerRemove kills and removes a container from the docker host.
func (m *ContainerAPIClientMock) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerRename changes the name of a given container.
func (m *ContainerAPIClientMock) ContainerRename(ctx context.Context, container, newContainerName string) error {
	args := m.Called(ctx, container, newContainerName)
	return args.Error(0)
}

// ContainerResize changes the size of the tty for a container.
func (m *ContainerAPIClientMock) ContainerResize(ctx context.Context, container string, options types.ResizeOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerRestart stops and starts a container again. It makes the daemon to
// wait for the container to be up again for a specific amount of time, given
// the timeout.
func (m *ContainerAPIClientMock) ContainerRestart(ctx context.Context, container string, timeout *time.Duration) error {
	args := m.Called(ctx, container, timeout)
	return args.Error(0)
}

// ContainerStatPath returns Stat information about a path inside the container
// filesystem.
func (m *ContainerAPIClientMock) ContainerStatPath(ctx context.Context, container, path string) (types.ContainerPathStat, error) {
	args := m.Called(ctx, container, path)
	return args.Get(0).(types.ContainerPathStat), args.Error(1)
}

// ContainerStats returns near realtime stats for a given container. It's up to
// the caller to close the io.ReadCloser returned.
func (m *ContainerAPIClientMock) ContainerStats(ctx context.Context, container string, stream bool) (types.ContainerStats, error) {
	args := m.Called(ctx, container, stream)
	return args.Get(0).(types.ContainerStats), args.Error(1)
}

// ContainerStart sends a request to the docker daemon to start a container.
func (m *ContainerAPIClientMock) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerStop stops a container without terminating the process. The process
// is blocked until the container stops or the timeout expires.
func (m *ContainerAPIClientMock) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	args := m.Called(ctx, container, timeout)
	return args.Error(0)
}

// ContainerTop shows process information from within a container.
func (m *ContainerAPIClientMock) ContainerTop(ctx context.Context, container string, arguments []string) (containertypes.ContainerTopOKBody, error) {
	args := m.Called(ctx, container, arguments)
	return args.Get(0).(containertypes.ContainerTopOKBody), args.Error(1)
}

// ContainerUnpause resumes the process execution within a container
func (m *ContainerAPIClientMock) ContainerUnpause(ctx context.Context, container string) error {
	args := m.Called(ctx, container)
	return args.Error(0)
}

// ContainerUpdate updates resources of a container.
func (m *ContainerAPIClientMock) ContainerUpdate(ctx context.Context, container string, updateConfig containertypes.UpdateConfig) (containertypes.ContainerUpdateOKBody, error) {
	args := m.Called(ctx, container, updateConfig)
	return args.Get(0).(containertypes.ContainerUpdateOKBody), args.Error(1)
}

// ContainerWait pauses execution until a container exits. It returns the API
// status code as response of its readiness.
func (m *ContainerAPIClientMock) ContainerWait(ctx context.Context, container string, condition containertypes.WaitCondition) (<-chan containertypes.ContainerWaitOKBody, <-chan error) {
	args := m.Called(ctx, container, condition)
	return args.Get(0).(<-chan containertypes.ContainerWaitOKBody), args.Get(1).(<-chan error)
}

// CopyFromContainer gets the content from the container and returns it as a
// Reader to manipulate it in the host. It's up to the caller to close the
// reader.
func (m *ContainerAPIClientMock) CopyFromContainer(ctx context.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error) {
	args := m.Called(ctx, container, srcPath)
	return args.Get(0).(io.ReadCloser), args.Get(1).(types.ContainerPathStat), args.Error(2)
}

// CopyToContainer copies content into the container filesystem.
func (m *ContainerAPIClientMock) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error {
	args := m.Called(ctx, container, path, content, options)
	return args.Error(0)
}

// ContainersPrune requests the daemon to delete unused data.
func (m *ContainerAPIClientMock) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(types.ContainersPruneReport), args.Error(1)
}

// DistributionAPIClientMock defines API client methods for the registry.
type DistributionAPIClientMock struct {
	*mock.Mock
}

// DistributionInspect returns the image digest with full Manifest.
func (m *DistributionAPIClientMock) DistributionInspect(ctx context.Context, image, encodedRegistryAuth string) (registry.DistributionInspect, error) {
	args := m.Called(ctx, image, encodedRegistryAuth)
	return args.Get(0).(registry.DistributionInspect), args.Error(1)
}

// ImageAPIClientMock defines API client methods for the images.
type ImageAPIClientMock struct {
	*mock.Mock
}

// ImageBuild sends request to the daemon to build images. The Body in the
// response implement an io.ReadCloser and it's up to the caller to close it.
func (m *ImageAPIClientMock) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	args := m.Called(ctx, context, options)
	return args.Get(0).(types.ImageBuildResponse), args.Error(1)
}

// BuildCachePrune requests the daemon to delete unused cache data.
func (m *ImageAPIClientMock) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.BuildCachePruneReport), args.Error(1)
}

// BuildCancel requests the daemon to cancel ongoing build request.
func (m *ImageAPIClientMock) BuildCancel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ImageCreate creates a new image based in the parent options. It returns the
// JSON content in the response body.
func (m *ImageAPIClientMock) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, parentReference, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageHistory returns the changes in an image in history format.
func (m *ImageAPIClientMock) ImageHistory(ctx context.Context, imageID string) ([]image.HistoryResponseItem, error) {
	args := m.Called(ctx, imageID)
	return args.Get(0).([]image.HistoryResponseItem), args.Error(1)
}

// ImageImport creates a new image based in the source options. It returns the
// JSON content in the response body.
func (m *ImageAPIClientMock) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, source, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageInspectWithRaw returns the image information and its raw
// representation.
func (m *ImageAPIClientMock) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	args := m.Called(ctx, image)
	return args.Get(0).(types.ImageInspect), args.Get(1).([]byte), args.Error(2)
}

// ImageList returns a list of images in the docker host.
func (m *ImageAPIClientMock) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.ImageSummary), args.Error(1)
}

// ImageLoad loads an image in the docker host from the client host. It's up to
// the caller to close the io.ReadCloser in the ImageLoadResponse returned by
// this function.
func (m *ImageAPIClientMock) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	args := m.Called(ctx, input, quiet)
	return args.Get(0).(types.ImageLoadResponse), args.Error(1)
}

// ImagePull requests the docker host to pull an image from a remote registry.
// It executes the privileged function if the operation is unauthorized and it
// tries one more time. It's up to the caller to handle the io.ReadCloser and
// close it properly.
func (m *ImageAPIClientMock) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImagePush requests the docker host to push an image to a remote registry. It
// executes the privileged function if the operation is unauthorized and it
// tries one more time. It's up to the caller to handle the io.ReadCloser and
// close it properly.
func (m *ImageAPIClientMock) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	args := m.Called(ctx)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageRemove removes an image from the docker host.
func (m *ImageAPIClientMock) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	args := m.Called(ctx, image, options)
	return args.Get(0).([]types.ImageDeleteResponseItem), args.Error(1)
}

// ImageSearch makes the docker host to search by a term in a remote registry.
// The list of results is not sorted in any fashion.
func (m *ImageAPIClientMock) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	args := m.Called(ctx, term, options)
	return args.Get(0).([]registry.SearchResult), args.Error(1)
}

// ImageSave retrieves one or more images from the docker host as an
// io.ReadCloser. It's up to the caller to store the images and close the
// stream.
func (m *ImageAPIClientMock) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	args := m.Called(ctx, images)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageTag tags an image in the docker host.
func (m *ImageAPIClientMock) ImageTag(ctx context.Context, image, ref string) error {
	args := m.Called(ctx, image, ref)
	return args.Error(0)
}

// ImagesPrune requests the daemon to delete unused data
func (m *ImageAPIClientMock) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.ImagesPruneReport), args.Error(0)
}

// NetworkAPIClientMock defines API client methods for the networks
type NetworkAPIClientMock struct {
	*mock.Mock
}

// NetworkConnect connects a container to an existent network in the docker
// host.
func (m *NetworkAPIClientMock) NetworkConnect(ctx context.Context, network, container string, config *networktypes.EndpointSettings) error {
	args := m.Called(ctx, network, container, config)
	return args.Error(0)
}

// NetworkCreate creates a new network in the docker host.
func (m *NetworkAPIClientMock) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(types.NetworkCreateResponse), args.Error(1)
}

// NetworkDisconnect disconnects a container from an existent network in the
// docker host.
func (m *NetworkAPIClientMock) NetworkDisconnect(ctx context.Context, network, container string, force bool) error {
	args := m.Called(ctx, network, container, force)
	return args.Error(0)
}

// NetworkInspect returns the information for a specific network configured in
// the docker host.
func (m *NetworkAPIClientMock) NetworkInspect(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, error) {
	args := m.Called(ctx, network, options)
	return args.Get(0).(types.NetworkResource), args.Error(1)
}

// NetworkInspectWithRaw returns the information for a specific network
// configured in the docker host and its raw representation.
func (m *NetworkAPIClientMock) NetworkInspectWithRaw(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, []byte, error) {
	args := m.Called(ctx, network, options)
	return args.Get(0).(types.NetworkResource), args.Get(1).([]byte), args.Error(2)
}

// NetworkList returns the list of networks configured in the docker host.
func (m *NetworkAPIClientMock) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.NetworkResource), args.Error(1)
}

// NetworkRemove removes an existent network from the docker host.
func (m *NetworkAPIClientMock) NetworkRemove(ctx context.Context, network string) error {
	args := m.Called(ctx, network)
	return args.Error(0)
}

// NetworksPrune requests the daemon to delete unused networks.
func (m *NetworkAPIClientMock) NetworksPrune(ctx context.Context, pruneFilter filters.Args) (types.NetworksPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.NetworksPruneReport), args.Error(1)
}

// NodeAPIClientMock defines API client methods for the nodes.
type NodeAPIClientMock struct {
	*mock.Mock
}

// NodeInspectWithRaw returns the node information.
func (m *NodeAPIClientMock) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error) {
	args := m.Called(ctx, nodeID)
	return args.Get(0).(swarm.Node), args.Get(1).([]byte), args.Error(2)
}

// NodeList returns the list of nodes.
func (m *NodeAPIClientMock) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Node), args.Error(1)
}

// NodeRemove removes a Node.
func (m *NodeAPIClientMock) NodeRemove(ctx context.Context, nodeID string, options types.NodeRemoveOptions) error {
	args := m.Called(ctx, nodeID, options)
	return args.Error(0)
}

// NodeUpdate updates a Node.
func (m *NodeAPIClientMock) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error {
	args := m.Called(ctx, nodeID, version, node)
	return args.Error(0)
}

// PluginAPIClientMock defines API client methods for the plugins.
type PluginAPIClientMock struct {
	*mock.Mock
}

// PluginList returns the installed plugins.
func (m *PluginAPIClientMock) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(types.PluginsListResponse), args.Error(1)
}

// PluginRemove removes a plugin.
func (m *PluginAPIClientMock) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginEnable enables a plugin.
func (m *PluginAPIClientMock) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginDisable disables a plugin.
func (m *PluginAPIClientMock) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginInstall installs a plugin.
func (m *PluginAPIClientMock) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginUpgrade upgrades a plugin.
func (m *PluginAPIClientMock) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginPush pushes a plugin to a registry.
func (m *PluginAPIClientMock) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error) {
	args := m.Called(ctx, name, registryAuth)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginSet modifies settings for an existing plugin.
func (m *PluginAPIClientMock) PluginSet(ctx context.Context, name string, args []string) error {
	arg := m.Called(ctx, name, args)
	return arg.Error(0)
}

// PluginInspectWithRaw inspects an existing plugin.
func (m *PluginAPIClientMock) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*types.Plugin), args.Get(1).([]byte), args.Error(2)
}

// PluginCreate creates a plugin.
func (m *PluginAPIClientMock) PluginCreate(ctx context.Context, createContext io.Reader, options types.PluginCreateOptions) error {
	args := m.Called(ctx, createContext, options)
	return args.Error(0)
}

// ServiceAPIClientMock defines API client methods for the services.
type ServiceAPIClientMock struct {
	*mock.Mock
}

// ServiceCreate creates a new Service.
func (m *ServiceAPIClientMock) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	args := m.Called(ctx, service, options)
	return args.Get(0).(types.ServiceCreateResponse), args.Error(0)
}

// ServiceInspectWithRaw returns the service information and the raw data.
func (m *ServiceAPIClientMock) ServiceInspectWithRaw(ctx context.Context, serviceID string, options types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	args := m.Called(ctx, serviceID, options)
	return args.Get(0).(swarm.Service), args.Get(1).([]byte), args.Error(2)
}

// ServiceList returns the list of services.
func (m *ServiceAPIClientMock) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Service), args.Error(1)
}

// ServiceRemove kills and removes a service.
func (m *ServiceAPIClientMock) ServiceRemove(ctx context.Context, serviceID string) error {
	args := m.Called(ctx, serviceID)
	return args.Error(0)
}

// ServiceUpdate updates a Service.
func (m *ServiceAPIClientMock) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	args := m.Called(ctx, serviceID, version, service, options)
	return args.Get(0).(types.ServiceUpdateResponse), args.Error(1)
}

// ServiceLogs returns the logs generated by a service in an io.ReadCloser.
// It's up to the caller to close the stream.
func (m *ServiceAPIClientMock) ServiceLogs(ctx context.Context, serviceID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, serviceID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// TaskLogs returns the logs generated by a task in an io.ReadCloser. It's up
// to the caller to close the stream.
func (m *ServiceAPIClientMock) TaskLogs(ctx context.Context, taskID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, taskID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// TaskInspectWithRaw returns the task information and its raw representation..
func (m *ServiceAPIClientMock) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(swarm.Task), args.Get(1).([]byte), args.Error(2)
}

// TaskList returns the list of tasks.
func (m *ServiceAPIClientMock) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Task), args.Error(1)
}

// SwarmAPIClientMock defines API client methods for the swarm.
type SwarmAPIClientMock struct {
	*mock.Mock
}

// SwarmInit initializes the Swarm.
func (m *SwarmAPIClientMock) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// SwarmJoin joins the Swarm.
func (m *SwarmAPIClientMock) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// SwarmGetUnlockKey retrieves the swarm's unlock key.
func (m *SwarmAPIClientMock) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.SwarmUnlockKeyResponse), args.Error(1)
}

// SwarmUnlock unlocks locked swarm.
func (m *SwarmAPIClientMock) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// SwarmLeave leaves the Swarm.
func (m *SwarmAPIClientMock) SwarmLeave(ctx context.Context, force bool) error {
	args := m.Called(ctx, force)
	return args.Error(0)
}

// SwarmInspect inspects the Swarm.
func (m *SwarmAPIClientMock) SwarmInspect(ctx context.Context) (swarm.Swarm, error) {
	args := m.Called(ctx)
	return args.Get(0).(swarm.Swarm), args.Error(1)
}

// SwarmUpdate updates the Swarm.
func (m *SwarmAPIClientMock) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error {
	args := m.Called(ctx, version, swarm, flags)
	return args.Error(0)
}

// SystemAPIClientMock defines API client methods for the system.
type SystemAPIClientMock struct {
	*mock.Mock
}

// Events returns a stream of events in the daemon. It's up to the caller to
// close the stream by cancelling the context. Once the stream has been
// completely read an io.EOF error will be sent over the error channel. If an
// error is sent all processing will be stopped. It's up to the caller to
// reopen the stream in the event of an error by reinvoking this method.
func (m *SystemAPIClientMock) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	args := m.Called(ctx, options)
	return args.Get(0).(<-chan events.Message), args.Get(1).(<-chan error)
}

// Info returns information about the docker server.
func (m *SystemAPIClientMock) Info(ctx context.Context) (types.Info, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Info), args.Error(1)
}

// RegistryLogin authenticates the docker server with a given docker registry.
// It returns UnauthorizerError when the authentication fails.
func (m *SystemAPIClientMock) RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	args := m.Called(ctx, auth)
	return args.Get(0).(registry.AuthenticateOKBody), args.Error(0)
}

// DiskUsage requests the current data usage from the daemon
func (m *SystemAPIClientMock) DiskUsage(ctx context.Context) (types.DiskUsage, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.DiskUsage), args.Error(1)
}

// Ping pings the server and return the value of the "Docker-Experimental"
// "API-Version" headers
func (m *SystemAPIClientMock) Ping(ctx context.Context) (types.Ping, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Ping), args.Error(1)
}

// VolumeAPIClientMock defines API client methods for the volumes.
type VolumeAPIClientMock struct {
	*mock.Mock
}

// VolumeCreate creates a volume in the docker host.
func (m *VolumeAPIClientMock) VolumeCreate(ctx context.Context, options volumetypes.VolumeCreateBody) (types.Volume, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(types.Volume), args.Error(1)
}

// VolumeInspect returns the information about a specific volume in the docker
// host.
func (m *VolumeAPIClientMock) VolumeInspect(ctx context.Context, volumeID string) (types.Volume, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(types.Volume), args.Error(1)
}

// VolumeInspectWithRaw returns the information about a specific volume in the
// docker host and its raw representation.
func (m *VolumeAPIClientMock) VolumeInspectWithRaw(ctx context.Context, volumeID string) (types.Volume, []byte, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(types.Volume), args.Get(1).([]byte), args.Error(2)
}

// VolumeList returns the volumes configured in the docker host.
func (m *VolumeAPIClientMock) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumeListOKBody, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(volumetypes.VolumeListOKBody), args.Error(1)
}

// VolumeRemove removes a volume from the docker host.
func (m *VolumeAPIClientMock) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	args := m.Called(ctx, volumeID, force)
	return args.Error(0)
}

// VolumesPrune requests the daemon to delete unused data.
func (m *VolumeAPIClientMock) VolumesPrune(ctx context.Context, pruneFilter filters.Args) (types.VolumesPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.VolumesPruneReport), args.Error(1)
}

// SecretAPIClientMock defines API client methods for secrets.
type SecretAPIClientMock struct {
	*mock.Mock
}

// SecretList returns the list of secrets.
func (m *SecretAPIClientMock) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Secret), args.Error(1)
}

// SecretCreate creates a new Secret.
func (m *SecretAPIClientMock) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error) {
	args := m.Called(ctx, secret)
	return args.Get(0).(types.SecretCreateResponse), args.Error(1)
}

// SecretRemove removes a Secret.
func (m *SecretAPIClientMock) SecretRemove(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SecretInspectWithRaw returns the secret information with raw data.
func (m *SecretAPIClientMock) SecretInspectWithRaw(ctx context.Context, name string) (swarm.Secret, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(swarm.Secret), args.Get(1).([]byte), args.Error(2)
}

// SecretUpdate updates a Secret. Currently, the only part of a secret spec
// which can be updated is Labels.
func (m *SecretAPIClientMock) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error {
	args := m.Called(ctx, id, version, secret)
	return args.Error(0)
}

// ConfigAPIClientMock defines API client methods for configs.
type ConfigAPIClientMock struct {
	*mock.Mock
}

// ConfigList returns the list of configs.
func (m *ConfigAPIClientMock) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Config), args.Error(1)
}

// ConfigCreate creates a new Config.
func (m *ConfigAPIClientMock) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(types.ConfigCreateResponse), args.Error(1)
}

// ConfigRemove removes a Config.
func (m *ConfigAPIClientMock) ConfigRemove(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ConfigInspectWithRaw returns the config information with raw data.
func (m *ConfigAPIClientMock) ConfigInspectWithRaw(ctx context.Context, name string) (swarm.Config, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(swarm.Config), args.Get(1).([]byte), args.Error(2)
}

// ConfigUpdate attempts to update a Config.
func (m *ConfigAPIClientMock) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error {
	args := m.Called(ctx, id, version, config)
	return args.Error(0)
}
