package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	var (
		dependencyKey = "1"
		s             = &Service{
			Hash: "1",
			Name: "TestStopRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.RUNNING, nil)
	mc.On("StopService", d.namespace(s.namespace())).Once().Return(nil)
	mc.On("DeleteNetwork", s.namespace()).Once().Return(nil)

	err := s.Stop(mc)
	require.NoError(t, err)

	mc.AssertExpectations(t)
}
