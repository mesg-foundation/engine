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

// ClientVersion mock.
func (m *CommonAPIClientMock) ClientVersion() string {
	args := m.Called()
	return args.String(0)
}

// DaemonHost mock.
func (m *CommonAPIClientMock) DaemonHost() string {
	args := m.Called()
	return args.String(0)
}

// HTTPClient mock.
func (m *CommonAPIClientMock) HTTPClient() *http.Client {
	args := m.Called()
	return args.Get(0).(*http.Client)
}

// ServerVersion mock.
func (m *CommonAPIClientMock) ServerVersion(ctx context.Context) (types.Version, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Version), args.Error(1)
}

// NegotiateAPIVersion mock.
func (m *CommonAPIClientMock) NegotiateAPIVersion(ctx context.Context) {
	m.Called(ctx)
}

// NegotiateAPIVersionPing mock.
func (m *CommonAPIClientMock) NegotiateAPIVersionPing(ping types.Ping) {
	m.Called(ping)
}

// DialSession mock.
func (m *CommonAPIClientMock) DialSession(ctx context.Context, proto string, meta map[string][]string) (net.Conn, error) {
	args := m.Called(ctx, proto, meta)
	return args.Get(0).(net.Conn), args.Error(1)
}

// Dialer mock.
func (m *CommonAPIClientMock) Dialer() func(context.Context) (net.Conn, error) {
	args := m.Called()
	return args.Get(0).(func(context.Context) (net.Conn, error))
}

// Close mock.
func (m *CommonAPIClientMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ContainerAPIClientMock defines API client methods for the containers
type ContainerAPIClientMock struct {
	*mock.Mock
}

// ContainerAttach mock.
func (m *ContainerAPIClientMock) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(types.HijackedResponse), args.Error(1)
}

// ContainerCommit mock.
func (m *ContainerAPIClientMock) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (types.IDResponse, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(types.IDResponse), args.Error(1)
}

// ContainerCreate mock.
func (m *ContainerAPIClientMock) ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, networkingConfig *networktypes.NetworkingConfig, containerName string) (containertypes.ContainerCreateCreatedBody, error) {
	args := m.Called(ctx, config, hostConfig, networkingConfig, containerName)
	return args.Get(0).(containertypes.ContainerCreateCreatedBody), args.Error(1)
}

// ContainerDiff mock.
func (m *ContainerAPIClientMock) ContainerDiff(ctx context.Context, container string) ([]containertypes.ContainerChangeResponseItem, error) {
	args := m.Called(ctx, container)
	return args.Get(0).([]containertypes.ContainerChangeResponseItem), args.Error(1)
}

// ContainerExecAttach mock.
func (m *ContainerAPIClientMock) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	args := m.Called(ctx, execID, config)
	return args.Get(0).(types.HijackedResponse), args.Error(1)
}

// ContainerExecCreate mock.
func (m *ContainerAPIClientMock) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	args := m.Called(ctx, container, config)
	return args.Get(0).(types.IDResponse), args.Error(1)
}

// ContainerExecInspect mock.
func (m *ContainerAPIClientMock) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	args := m.Called(ctx, execID)
	return args.Get(0).(types.ContainerExecInspect), args.Error(1)
}

// ContainerExecResize mock.
func (m *ContainerAPIClientMock) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	args := m.Called(ctx, execID, options)
	return args.Error(0)
}

// ContainerExecStart mock.
func (m *ContainerAPIClientMock) ContainerExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	args := m.Called(ctx, execID, config)
	return args.Error(0)
}

// ContainerExport mock.
func (m *ContainerAPIClientMock) ContainerExport(ctx context.Context, container string) (io.ReadCloser, error) {
	args := m.Called(ctx, container)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ContainerInspect mock.
func (m *ContainerAPIClientMock) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	args := m.Called(ctx, container)
	return args.Get(0).(types.ContainerJSON), args.Error(1)
}

// ContainerInspectWithRaw mock.
func (m *ContainerAPIClientMock) ContainerInspectWithRaw(ctx context.Context, container string, getSize bool) (types.ContainerJSON, []byte, error) {
	args := m.Called(ctx, container, getSize)
	return args.Get(0).(types.ContainerJSON), args.Get(1).([]byte), args.Error(2)
}

// ContainerKill mock.
func (m *ContainerAPIClientMock) ContainerKill(ctx context.Context, container, signal string) error {
	args := m.Called(ctx, container, signal)
	return args.Error(0)
}

// ContainerList mock.
func (m *ContainerAPIClientMock) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}

// ContainerLogs mock.
func (m *ContainerAPIClientMock) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ContainerPause mock.
func (m *ContainerAPIClientMock) ContainerPause(ctx context.Context, container string) error {
	args := m.Called(ctx, container)
	return args.Error(0)
}

// ContainerRemove mock.
func (m *ContainerAPIClientMock) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerRename mock.
func (m *ContainerAPIClientMock) ContainerRename(ctx context.Context, container, newContainerName string) error {
	args := m.Called(ctx, container, newContainerName)
	return args.Error(0)
}

// ContainerResize mock.
func (m *ContainerAPIClientMock) ContainerResize(ctx context.Context, container string, options types.ResizeOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerRestart mock.
func (m *ContainerAPIClientMock) ContainerRestart(ctx context.Context, container string, timeout *time.Duration) error {
	args := m.Called(ctx, container, timeout)
	return args.Error(0)
}

// ContainerStatPath mock.
func (m *ContainerAPIClientMock) ContainerStatPath(ctx context.Context, container, path string) (types.ContainerPathStat, error) {
	args := m.Called(ctx, container, path)
	return args.Get(0).(types.ContainerPathStat), args.Error(1)
}

// ContainerStats mock.
func (m *ContainerAPIClientMock) ContainerStats(ctx context.Context, container string, stream bool) (types.ContainerStats, error) {
	args := m.Called(ctx, container, stream)
	return args.Get(0).(types.ContainerStats), args.Error(1)
}

// ContainerStart mock.
func (m *ContainerAPIClientMock) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	args := m.Called(ctx, container, options)
	return args.Error(0)
}

// ContainerStop mock.
func (m *ContainerAPIClientMock) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	args := m.Called(ctx, container, timeout)
	return args.Error(0)
}

// ContainerTop mock.
func (m *ContainerAPIClientMock) ContainerTop(ctx context.Context, container string, arguments []string) (containertypes.ContainerTopOKBody, error) {
	args := m.Called(ctx, container, arguments)
	return args.Get(0).(containertypes.ContainerTopOKBody), args.Error(1)
}

// ContainerUnpause mock.
func (m *ContainerAPIClientMock) ContainerUnpause(ctx context.Context, container string) error {
	args := m.Called(ctx, container)
	return args.Error(0)
}

// ContainerUpdate mock.
func (m *ContainerAPIClientMock) ContainerUpdate(ctx context.Context, container string, updateConfig containertypes.UpdateConfig) (containertypes.ContainerUpdateOKBody, error) {
	args := m.Called(ctx, container, updateConfig)
	return args.Get(0).(containertypes.ContainerUpdateOKBody), args.Error(1)
}

// ContainerWait mock.
func (m *ContainerAPIClientMock) ContainerWait(ctx context.Context, container string, condition containertypes.WaitCondition) (<-chan containertypes.ContainerWaitOKBody, <-chan error) {
	args := m.Called(ctx, container, condition)
	return args.Get(0).(<-chan containertypes.ContainerWaitOKBody), args.Get(1).(<-chan error)
}

// CopyFromContainer mock.
func (m *ContainerAPIClientMock) CopyFromContainer(ctx context.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error) {
	args := m.Called(ctx, container, srcPath)
	return args.Get(0).(io.ReadCloser), args.Get(1).(types.ContainerPathStat), args.Error(2)
}

// CopyToContainer mock.
func (m *ContainerAPIClientMock) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error {
	args := m.Called(ctx, container, path, content, options)
	return args.Error(0)
}

// ContainersPrune mock.
func (m *ContainerAPIClientMock) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error) {
	args := m.Called(ctx, pruneFilters)
	return args.Get(0).(types.ContainersPruneReport), args.Error(1)
}

// DistributionAPIClientMock defines API client methods for the registry
type DistributionAPIClientMock struct {
	*mock.Mock
}

// DistributionInspect mock.
func (m *DistributionAPIClientMock) DistributionInspect(ctx context.Context, image, encodedRegistryAuth string) (registry.DistributionInspect, error) {
	args := m.Called(ctx, image, encodedRegistryAuth)
	return args.Get(0).(registry.DistributionInspect), args.Error(1)
}

// ImageAPIClientMock defines API client methods for the images
type ImageAPIClientMock struct {
	*mock.Mock
}

// ImageBuild mock.
func (m *ImageAPIClientMock) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	args := m.Called(ctx, context, options)
	return args.Get(0).(types.ImageBuildResponse), args.Error(1)
}

// BuildCachePrune mock.
func (m *ImageAPIClientMock) BuildCachePrune(ctx context.Context) (*types.BuildCachePruneReport, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.BuildCachePruneReport), args.Error(1)
}

// BuildCancel mock.
func (m *ImageAPIClientMock) BuildCancel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ImageCreate mock.
func (m *ImageAPIClientMock) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, parentReference, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageHistory mock.
func (m *ImageAPIClientMock) ImageHistory(ctx context.Context, imageID string) ([]image.HistoryResponseItem, error) {
	args := m.Called(ctx, imageID)
	return args.Get(0).([]image.HistoryResponseItem), args.Error(1)
}

// ImageImport mock.
func (m *ImageAPIClientMock) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, source, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageInspectWithRaw mock.
func (m *ImageAPIClientMock) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	args := m.Called(ctx, image)
	return args.Get(0).(types.ImageInspect), args.Get(1).([]byte), args.Error(2)
}

// ImageList mock.
func (m *ImageAPIClientMock) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.ImageSummary), args.Error(1)
}

// ImageLoad mock.
func (m *ImageAPIClientMock) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	args := m.Called(ctx, input, quiet)
	return args.Get(0).(types.ImageLoadResponse), args.Error(1)
}

// ImagePull mock.
func (m *ImageAPIClientMock) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImagePush mock.
func (m *ImageAPIClientMock) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	args := m.Called(ctx)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageRemove mock.
func (m *ImageAPIClientMock) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	args := m.Called(ctx, image, options)
	return args.Get(0).([]types.ImageDeleteResponseItem), args.Error(1)
}

// ImageSearch mock.
func (m *ImageAPIClientMock) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	args := m.Called(ctx, term, options)
	return args.Get(0).([]registry.SearchResult), args.Error(1)
}

// ImageSave mock.
func (m *ImageAPIClientMock) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	args := m.Called(ctx, images)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageTag mock.
func (m *ImageAPIClientMock) ImageTag(ctx context.Context, image, ref string) error {
	args := m.Called(ctx, image, ref)
	return args.Error(0)
}

// ImagesPrune mock.
func (m *ImageAPIClientMock) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.ImagesPruneReport), args.Error(0)
}

// NetworkAPIClientMock defines API client methods for the networks
type NetworkAPIClientMock struct {
	*mock.Mock
}

// NetworkConnect mock.
func (m *NetworkAPIClientMock) NetworkConnect(ctx context.Context, network, container string, config *networktypes.EndpointSettings) error {
	args := m.Called(ctx, network, container, config)
	return args.Error(0)
}

// NetworkCreate mock.
func (m *NetworkAPIClientMock) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(types.NetworkCreateResponse), args.Error(1)
}

// NetworkDisconnect mock.
func (m *NetworkAPIClientMock) NetworkDisconnect(ctx context.Context, network, container string, force bool) error {
	args := m.Called(ctx, network, container, force)
	return args.Error(0)
}

// NetworkInspect mock.
func (m *NetworkAPIClientMock) NetworkInspect(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.NetworkResource), args.Error(1)
}

// NetworkInspectWithRaw mock.
func (m *NetworkAPIClientMock) NetworkInspectWithRaw(ctx context.Context, network string, options types.NetworkInspectOptions) (types.NetworkResource, []byte, error) {
	args := m.Called(ctx, network, options)
	return args.Get(0).(types.NetworkResource), args.Get(1).([]byte), args.Error(2)
}

// NetworkList mock.
func (m *NetworkAPIClientMock) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.NetworkResource), args.Error(1)
}

// NetworkRemove mock.
func (m *NetworkAPIClientMock) NetworkRemove(ctx context.Context, network string) error {
	args := m.Called(ctx, network)
	return args.Error(0)
}

// NetworksPrune mock.
func (m *NetworkAPIClientMock) NetworksPrune(ctx context.Context, pruneFilter filters.Args) (types.NetworksPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.NetworksPruneReport), args.Error(1)
}

// NodeAPIClientMock defines API client methods for the nodes
type NodeAPIClientMock struct {
	*mock.Mock
}

// NodeInspectWithRaw mock.
func (m *NodeAPIClientMock) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error) {
	args := m.Called(ctx, nodeID)
	return args.Get(0).(swarm.Node), args.Get(1).([]byte), args.Error(2)
}

// NodeList mock.
func (m *NodeAPIClientMock) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Node), args.Error(1)
}

// NodeRemove mock.
func (m *NodeAPIClientMock) NodeRemove(ctx context.Context, nodeID string, options types.NodeRemoveOptions) error {
	args := m.Called(ctx, nodeID, options)
	return args.Error(0)
}

// NodeUpdate mock.
func (m *NodeAPIClientMock) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error {
	args := m.Called(ctx, nodeID, version, node)
	return args.Error(0)
}

// PluginAPIClientMock defines API client methods for the plugins
type PluginAPIClientMock struct {
	*mock.Mock
}

// PluginList mock.
func (m *PluginAPIClientMock) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(types.PluginsListResponse), args.Error(1)
}

// PluginRemove mock.
func (m *PluginAPIClientMock) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginEnable mock.
func (m *PluginAPIClientMock) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginDisable mock.
func (m *PluginAPIClientMock) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error {
	args := m.Called(ctx, name, options)
	return args.Error(0)
}

// PluginInstall mock.
func (m *PluginAPIClientMock) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginUpgrade mock.
func (m *PluginAPIClientMock) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, name, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginPush mock.
func (m *PluginAPIClientMock) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error) {
	args := m.Called(ctx, name, registryAuth)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// PluginSet mock.
func (m *PluginAPIClientMock) PluginSet(ctx context.Context, name string, args []string) error {
	arg := m.Called(ctx, name, args)
	return arg.Error(0)
}

// PluginInspectWithRaw mock.
func (m *PluginAPIClientMock) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*types.Plugin), args.Get(1).([]byte), args.Error(2)
}

// PluginCreate mock.
func (m *PluginAPIClientMock) PluginCreate(ctx context.Context, createContext io.Reader, options types.PluginCreateOptions) error {
	args := m.Called(ctx, createContext, options)
	return args.Error(0)
}

// ServiceAPIClientMock defines API client methods for the services
type ServiceAPIClientMock struct {
	*mock.Mock
}

// ServiceCreate mock.
func (m *ServiceAPIClientMock) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	args := m.Called(ctx, service, options)
	return args.Get(0).(types.ServiceCreateResponse), args.Error(0)
}

// ServiceInspectWithRaw mock.
func (m *ServiceAPIClientMock) ServiceInspectWithRaw(ctx context.Context, serviceID string, options types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	args := m.Called(ctx, serviceID, options)
	return args.Get(0).(swarm.Service), args.Get(1).([]byte), args.Error(2)
}

// ServiceList mock.
func (m *ServiceAPIClientMock) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Service), args.Error(1)
}

// ServiceRemove mock.
func (m *ServiceAPIClientMock) ServiceRemove(ctx context.Context, serviceID string) error {
	args := m.Called(ctx, serviceID)
	return args.Error(0)
}

// ServiceUpdate mock.
func (m *ServiceAPIClientMock) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	args := m.Called(ctx, serviceID, version, service, options)
	return args.Get(0).(types.ServiceUpdateResponse), args.Error(1)
}

// ServiceLogs mock.
func (m *ServiceAPIClientMock) ServiceLogs(ctx context.Context, serviceID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, serviceID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// TaskLogs mock.
func (m *ServiceAPIClientMock) TaskLogs(ctx context.Context, taskID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, taskID, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// TaskInspectWithRaw mock.
func (m *ServiceAPIClientMock) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(swarm.Task), args.Get(1).([]byte), args.Error(2)
}

// TaskList mock.
func (m *ServiceAPIClientMock) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Task), args.Error(1)
}

// SwarmAPIClientMock defines API client methods for the swarm
type SwarmAPIClientMock struct {
	*mock.Mock
}

// SwarmInit mock.
func (m *SwarmAPIClientMock) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// SwarmJoin mock.
func (m *SwarmAPIClientMock) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// SwarmGetUnlockKey mock.
func (m *SwarmAPIClientMock) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.SwarmUnlockKeyResponse), args.Error(1)
}

// SwarmUnlock mock.
func (m *SwarmAPIClientMock) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// SwarmLeave mock.
func (m *SwarmAPIClientMock) SwarmLeave(ctx context.Context, force bool) error {
	args := m.Called(ctx, force)
	return args.Error(0)
}

// SwarmInspect mock.
func (m *SwarmAPIClientMock) SwarmInspect(ctx context.Context) (swarm.Swarm, error) {
	args := m.Called(ctx)
	return args.Get(0).(swarm.Swarm), args.Error(1)
}

// SwarmUpdate mock.
func (m *SwarmAPIClientMock) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, flags swarm.UpdateFlags) error {
	args := m.Called(ctx, version, swarm, flags)
	return args.Error(0)
}

// SystemAPIClientMock defines API client methods for the system
type SystemAPIClientMock struct {
	*mock.Mock
}

// Events mock.
func (m *SystemAPIClientMock) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	args := m.Called(ctx, options)
	return args.Get(0).(<-chan events.Message), args.Get(1).(<-chan error)
}

// Info mock.
func (m *SystemAPIClientMock) Info(ctx context.Context) (types.Info, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Info), args.Error(1)
}

// RegistryLogin mock.
func (m *SystemAPIClientMock) RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	args := m.Called(ctx, auth)
	return args.Get(0).(registry.AuthenticateOKBody), args.Error(0)
}

// DiskUsage mock.
func (m *SystemAPIClientMock) DiskUsage(ctx context.Context) (types.DiskUsage, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.DiskUsage), args.Error(1)
}

// Ping mock.
func (m *SystemAPIClientMock) Ping(ctx context.Context) (types.Ping, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.Ping), args.Error(1)
}

// VolumeAPIClientMock defines API client methods for the volumes
type VolumeAPIClientMock struct {
	*mock.Mock
}

// VolumeCreate mock.
func (m *VolumeAPIClientMock) VolumeCreate(ctx context.Context, options volumetypes.VolumeCreateBody) (types.Volume, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(types.Volume), args.Error(1)
}

// VolumeInspect mock.
func (m *VolumeAPIClientMock) VolumeInspect(ctx context.Context, volumeID string) (types.Volume, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(types.Volume), args.Error(1)
}

// VolumeInspectWithRaw mock.
func (m *VolumeAPIClientMock) VolumeInspectWithRaw(ctx context.Context, volumeID string) (types.Volume, []byte, error) {
	args := m.Called(ctx, volumeID)
	return args.Get(0).(types.Volume), args.Get(1).([]byte), args.Error(2)
}

// VolumeList mock.
func (m *VolumeAPIClientMock) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumeListOKBody, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(volumetypes.VolumeListOKBody), args.Error(1)
}

// VolumeRemove mock.
func (m *VolumeAPIClientMock) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	args := m.Called(ctx, volumeID, force)
	return args.Error(0)
}

// VolumesPrune mock.
func (m *VolumeAPIClientMock) VolumesPrune(ctx context.Context, pruneFilter filters.Args) (types.VolumesPruneReport, error) {
	args := m.Called(ctx, pruneFilter)
	return args.Get(0).(types.VolumesPruneReport), args.Error(1)
}

// SecretAPIClientMock defines API client methods for secrets
type SecretAPIClientMock struct {
	*mock.Mock
}

// SecretList mock.
func (m *SecretAPIClientMock) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Secret), args.Error(1)
}

// SecretCreate mock.
func (m *SecretAPIClientMock) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error) {
	args := m.Called(ctx, secret)
	return args.Get(0).(types.SecretCreateResponse), args.Error(1)
}

// SecretRemove mock.
func (m *SecretAPIClientMock) SecretRemove(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SecretInspectWithRaw mock.
func (m *SecretAPIClientMock) SecretInspectWithRaw(ctx context.Context, name string) (swarm.Secret, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(swarm.Secret), args.Get(1).([]byte), args.Error(2)
}

// SecretUpdate mock.
func (m *SecretAPIClientMock) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error {
	args := m.Called(ctx, id, version, secret)
	return args.Error(0)
}

// ConfigAPIClientMock defines API client methods for configs
type ConfigAPIClientMock struct {
	*mock.Mock
}

// ConfigList mock.
func (m *ConfigAPIClientMock) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]swarm.Config), args.Error(1)
}

// ConfigCreate mock.
func (m *ConfigAPIClientMock) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(types.ConfigCreateResponse), args.Error(1)
}

// ConfigRemove mock.
func (m *ConfigAPIClientMock) ConfigRemove(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ConfigInspectWithRaw mock.
func (m *ConfigAPIClientMock) ConfigInspectWithRaw(ctx context.Context, name string) (swarm.Config, []byte, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(swarm.Config), args.Get(1).([]byte), args.Error(2)
}

// ConfigUpdate mock.
func (m *ConfigAPIClientMock) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error {
	args := m.Called(ctx, id, version, config)
	return args.Error(0)
}
