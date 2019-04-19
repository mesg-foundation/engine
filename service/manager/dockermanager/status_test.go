package dockermanager

import (
	"errors"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestUnknownServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		statusErr     = errors.New("ops")
		s             = &service.Service{
			Hash: "1",
			Name: "TestUnknownServiceStatus",
			Dependencies: []*service.Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
		m  = New(mc)
	)

	d, _ := s.GetDependency(dependencyKey)

	mc.On("Status", d.Namespace(s.Namespace())).Once().Return(container.UNKNOWN, statusErr)

	status, err := m.Status(s)
	require.Equal(t, statusErr, err)
	require.Equal(t, service.UNKNOWN, status)

	mc.AssertExpectations(t)
}

func TestStoppedServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		s             = &service.Service{
			Hash: "1",
			Name: "TestStoppedServiceStatus",
			Dependencies: []*service.Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
		m  = New(mc)
	)

	d, _ := s.GetDependency(dependencyKey)

	mc.On("Status", d.Namespace(s.Namespace())).Once().Return(container.STOPPED, nil)

	status, err := m.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.STOPPED, status)

	mc.AssertExpectations(t)
}

func TestRunningServiceStatus(t *testing.T) {
	var (
		dependencyKey = "1"
		s             = &service.Service{
			Hash: "1",
			Name: "TestRunningServiceStatus",
			Dependencies: []*service.Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
		m  = New(mc)
	)

	d, _ := s.GetDependency(dependencyKey)

	mc.On("Status", d.Namespace(s.Namespace())).Once().Return(container.RUNNING, nil)

	status, err := m.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.RUNNING, status)

	mc.AssertExpectations(t)
}

func TestPartialServiceStatus(t *testing.T) {
	var (
		dependencyKey  = "1"
		dependencyKey2 = "2"
		s              = &service.Service{
			Hash: "1",
			Name: "TestPartialServiceStatus",
			Dependencies: []*service.Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
				{
					Key:   dependencyKey2,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
		m  = New(mc)
	)

	var (
		d, _  = s.GetDependency(dependencyKey)
		d2, _ = s.GetDependency(dependencyKey2)
	)

	mc.On("Status", d.Namespace(s.Namespace())).Once().Return(container.RUNNING, nil)
	mc.On("Status", d2.Namespace(s.Namespace())).Once().Return(container.STOPPED, nil)

	status, err := m.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.PARTIAL, status)

	mc.AssertExpectations(t)
}
