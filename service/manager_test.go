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
		s = &Service{
			Configuration: Dependency{
				Env: []string{"foo="},
			},
		}
		c   = &mocks.Container{}
		m   = NewContainerManager(c, nil)
		env = map[string]string{"foo": "bar"}
	)
	c.On("Build", "testdata").Return("ff", nil)
	assert.NoError(t, m.Deploy(s, "testdata", env))

	assert.Equal(t, "ff", s.Configuration.Image)
	assert.Equal(t, []string{"foo=bar"}, s.Configuration.Env)
	c.AssertExpectations(t)
}

func TestContainerManagerStart(t *testing.T) {
	var cfg config.Config
	cfg.Server.Address = ":80"
	cfg.Core.Name = "core"

	var (
		s = &Service{
			Dependencies: map[string]*Dependency{
				"dummy": {},
			},
		}
		c = &mocks.Container{}
		m = NewContainerManager(c, &cfg)
	)

	c.On("Status", s.namespace()).Return(container.STOPPED, nil)
	c.On("Status", depNamespace(s.Hash, "dummy")).Return(container.STOPPED, nil)

	c.On("CreateNetwork", s.namespace()).Return("network", nil)
	c.On("SharedNetworkID").Return("shared-network", nil)
	c.On("StartService", mock.Anything).Return("", nil)
	c.On("StartService", mock.Anything).Return("", nil)
	assert.NoError(t, m.Start(s))

	assert.Equal(t, StatusRunning, s.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerStop(t *testing.T) {
	var (
		s = &Service{
			Dependencies: map[string]*Dependency{
				"dummy": nil,
			},
		}
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("StopService", s.namespace()).Return(nil)
	c.On("StopService", depNamespace(s.Hash, "dummy")).Return(nil)
	c.On("DeleteNetwork", s.namespace(), container.EventDestroy).Return(nil)
	assert.NoError(t, m.Stop(s))

	assert.Equal(t, StatusStopped, s.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerDelete(t *testing.T) {
	var (
		s = &Service{
			Configuration: Dependency{
				Volumes: []string{"foo"},
			},
			Dependencies: map[string]*Dependency{
				"dummy": {
					Volumes: []string{"bar"},
				},
			},
		}
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("DeleteVolume", s.volumes(MainServiceKey)[0].Source).Return(nil)
	c.On("DeleteVolume", s.volumes("dummy")[0].Source).Return(nil)
	assert.NoError(t, m.Delete(s))

	assert.Equal(t, StatusDeleted, s.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerStatus(t *testing.T) {
	var (
		s = &Service{
			Dependencies: map[string]*Dependency{
				"dummy": nil,
			},
		}
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)
	)

	c.On("Status", s.namespace()).Return(container.RUNNING, nil)
	c.On("Status", depNamespace(s.Hash, "dummy")).Return(container.RUNNING, nil)
	assert.NoError(t, m.Status(s))

	assert.Equal(t, StatusRunning, s.Status)
	c.AssertExpectations(t)
}

func TestContainerManagerLogs(t *testing.T) {
	var (
		s = &Service{
			Dependencies: map[string]*Dependency{
				"dummy": nil,
			},
		}
		c = &mocks.Container{}
		m = NewContainerManager(c, nil)

		l1 = ioutil.NopCloser(nil)
		l2 = ioutil.NopCloser(nil)
	)

	c.On("ServiceLogs", s.namespace()).Return(l1, nil)
	c.On("ServiceLogs", depNamespace(s.Hash, "dummy")).Return(l2, nil)
	logs, err := m.Logs(s, nil)
	assert.NoError(t, err)

	assert.Equal(t, MainServiceKey, logs[0].key)
	assert.Equal(t, l1, logs[0].r)
	assert.Equal(t, "dummy", logs[1].key)
	assert.Equal(t, l2, logs[1].r)
	c.AssertExpectations(t)
}
