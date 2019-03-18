package service

import (
	"io/ioutil"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContainerManagerDeploy(t *testing.T) {
	var (
		c   = &mocks.Container{}
		m   = NewContainerManager(c, nil)
		env = map[string]string{"foo": "bar"}
	)
	c.On("Build", "testdata").Return("ff", nil)
	assert.NoError(t, m.Deploy(ts, "testdata", env))

	assert.Equal(t, "ff", ts.Configuration.Image)
	assert.Equal(t, []string{"foo=bar"}, ts.Configuration.Env)
	c.AssertExpectations(t)
}

func TestContainerManagerStart(t *testing.T) {
	var cfg config.Config
	cfg.Server.Address = ":80"
	cfg.Core.Name = "core"

	var (
		c = &mocks.Container{}
		m = NewContainerManager(c, &cfg)
	)

	c.On("Status", ts.namespace()).Return(container.STOPPED, nil)
	c.On("Status", depNamespace(ts.Hash, "dummy")).Return(container.STOPPED, nil)

	c.On("CreateNetwork", ts.namespace()).Return("network", nil)
	c.On("SharedNetworkID").Return("shared-network", nil)
	c.On("StartService", mock.Anything).Return("", nil)
	c.On("StartService", mock.Anything).Return("", nil)
	assert.NoError(t, m.Start(ts))

	assert.Equal(t, StatusRunning, ts.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerStop(t *testing.T) {
	var (
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("StopService", ts.namespace()).Return(nil)
	c.On("StopService", depNamespace(ts.Hash, "dummy")).Return(nil)
	c.On("DeleteNetwork", ts.namespace(), container.EventDestroy).Return(nil)
	assert.NoError(t, m.Stop(ts))

	assert.Equal(t, StatusStopped, ts.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerDelete(t *testing.T) {
	var (
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("DeleteVolume", ts.volumes(MainServiceKey)[0].Source).Return(nil)
	c.On("DeleteVolume", ts.volumes("dummy")[0].Source).Return(nil)
	assert.NoError(t, m.Delete(ts))

	assert.Equal(t, StatusDeleted, ts.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerStatus(t *testing.T) {
	var (
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("Status", ts.namespace()).Return(container.RUNNING, nil)
	c.On("Status", depNamespace(ts.Hash, "dummy")).Return(container.RUNNING, nil)
	assert.NoError(t, m.Status(ts))

	assert.Equal(t, StatusRunning, ts.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerLogs(t *testing.T) {
	var (
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)

		l1 = ioutil.NopCloser(nil)
		l2 = ioutil.NopCloser(nil)
	)

	c.On("ServiceLogs", ts.namespace()).Return(l1, nil)
	c.On("ServiceLogs", depNamespace(ts.Hash, "dummy")).Return(l2, nil)
	logs, err := m.Logs(ts, nil)
	assert.NoError(t, err)

	assert.Equal(t, MainServiceKey, logs[0].key)
	assert.Equal(t, l1, logs[0].r)
	assert.Equal(t, "dummy", logs[1].key)
	assert.Equal(t, l2, logs[1].r)
	c.AssertExpectations(t)
}
