package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	var (
		dependencyKey = "1"
		s             = &service.Service{
			Hash: []byte{0},
			Name: "TestStopRunningService",
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

	mc.On("Status", dependencyNamespace(serviceNamespace(s.Hash), d.Key)).Once().Return(container.RUNNING, nil)
	mc.On("StopService", dependencyNamespace(serviceNamespace(s.Hash), d.Key)).Once().Return(nil)
	mc.On("DeleteNetwork", serviceNamespace(s.Hash)).Once().Return(nil)

	err := m.Stop(s)
	require.NoError(t, err)

	mc.AssertExpectations(t)
}
