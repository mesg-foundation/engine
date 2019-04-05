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

	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace(s.namespace())).Once().Return(nil)
	mc.On("DeleteNetwork", s.namespace()).Once().Return(nil)

	err := s.Stop(mc)
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
	mc.On("StopService", d.namespace(s.namespace())).Once().Return(nil)
	require.NoError(t, d.Stop(mc, s))
	mc.AssertExpectations(t)
}
