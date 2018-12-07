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

	mc.On("Status", d.namespace()).Once().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace()).Once().Return(nil)
	mc.On("DeleteNetwork", s.namespace(), container.EventDestroy).Once().Return(nil)

	err := s.Stop()
	require.NoError(t, err)

	mc.AssertExpectations(t)
}

func TestStopDependency(t *testing.T) {
	var (
		dependencyKey = "1"
		s, mc         = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestStopService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	d, _ := s.getDependency(dependencyKey)
	mc.On("StopService", d.namespace()).Once().Return(nil)
	require.NoError(t, d.Stop())
	mc.AssertExpectations(t)
}
