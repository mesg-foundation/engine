package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	cfg, _ := config.Global()
	c := &mocks.Container{}
	d := NewContainerDaemon(cfg, c)

	c.On("SharedNetworkID").Return("1", nil)
	c.On("StartService", d.buildServiceOptions("1")).Return("1", nil)
	require.NoError(t, d.Start())
	c.AssertExpectations(t)
}

func TestStop(t *testing.T) {
	cfg, _ := config.Global()
	c := &mocks.Container{}
	d := NewContainerDaemon(cfg, c)

	c.On("StopService", []string{}).Return(nil)
	require.NoError(t, d.Stop())
	c.AssertExpectations(t)
}

func TestStatus(t *testing.T) {
	cfg, _ := config.Global()
	c := &mocks.Container{}
	d := NewContainerDaemon(cfg, c)

	c.On("Status", []string{}).Return(container.STOPPED, nil)
	status, err := d.Status()
	require.NoError(t, err)
	require.Equal(t, container.STOPPED, status)
	c.AssertExpectations(t)
}

func TestLogs(t *testing.T) {
	cfg, _ := config.Global()
	c := &mocks.Container{}
	d := NewContainerDaemon(cfg, c)

	c.On("ServiceLogs", []string{}).Return(nil, nil)
	_, err := d.Logs()
	require.NoError(t, err)
	c.AssertExpectations(t)
}

func TestBuildServiceOptions(t *testing.T) {
	cfg, _ := config.Global()
	c := &mocks.Container{}
	d := NewContainerDaemon(cfg, c)

	spec := d.buildServiceOptions("")
	require.Equal(t, []string{}, spec.Namespace)
	// Make sure that the config directory is passed in parameter to write on the same folder
	require.Contains(t, spec.Env, "MESG_LOG_LEVEL=info")
	require.Contains(t, spec.Env, "MESG_LOG_FORMAT=text")
	require.Contains(t, spec.Env, "MESG_CORE_PATH="+cfg.Docker.Core.Path)
	// Ensure that the port is shared
	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)
	require.Equal(t, spec.Ports[0].Published, uint32(port))
	require.Equal(t, spec.Ports[0].Target, uint32(port))
	// Ensure that the docker socket is shared in the core
	require.Equal(t, spec.Mounts[0].Source, cfg.Docker.Socket)
	require.Equal(t, spec.Mounts[0].Target, cfg.Docker.Socket)
	require.True(t, spec.Mounts[0].Bind)
	// Ensure that the host users folder is sync with the core
	require.Equal(t, spec.Mounts[1].Source, cfg.Core.Path)
	require.Equal(t, spec.Mounts[1].Target, cfg.Docker.Core.Path)
	require.True(t, spec.Mounts[1].Bind)
}
