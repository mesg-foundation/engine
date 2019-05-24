package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestDeleteVolumes(t *testing.T) {
	var (
		dependencyKey1 = "1"
		dependencyKey2 = "2"
		volumeA        = "a"
		volumeB        = "b"
		s              = &service.Service{
			Name: "TestCreateVolumes",
			Dependencies: []*service.Dependency{
				{
					Key:     dependencyKey1,
					Image:   "1",
					Volumes: []string{volumeA, volumeB},
				},
				{
					Key:         dependencyKey2,
					Image:       "1",
					VolumesFrom: []string{dependencyKey1},
				},
			},
		}
		mc = &mocks.Container{}
		m  = New(mc)
	)

	var (
		d1, _    = s.GetDependency(dependencyKey1)
		volumes1 = extractVolumes(s, d1)
	)

	mc.On("DeleteVolume", volumes1[0].Source).Once().Return(nil)
	mc.On("DeleteVolume", volumes1[1].Source).Once().Return(nil)

	require.NoError(t, m.Delete(s))

	mc.AssertExpectations(t)
}
