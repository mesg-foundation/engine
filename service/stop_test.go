package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Twice().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace()).Once().Return(nil)
	mc.On("DeleteNetwork", s.namespace()).Once().Return(nil)

	err := s.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}

func TestStopNonRunningService(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopNonRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.STOPPED, nil)

	err := s.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}

func TestStopRunningDependency(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopNonRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace()).Once().Return(nil)

	err := d.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}

func TestStopNonRunningDependency(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopNonRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace()).Once().Return(container.STOPPED, nil)

	err := d.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}
