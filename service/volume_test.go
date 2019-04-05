package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteVolumes(t *testing.T) {
	var (
		dependencyKey1 = "1"
		dependencyKey2 = "2"
		volumeA        = "a"
		volumeB        = "b"
		s, mc          = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestCreateVolumes",
			Dependencies: []*Dependency{
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
		})
	)

	var (
		d1, _    = s.getDependency(dependencyKey1)
		volumes1 = d1.extractVolumes(s)
	)

	mc.On("DeleteVolume", volumes1[0].Source).Once().Return(nil)
	mc.On("DeleteVolume", volumes1[1].Source).Once().Return(nil)

	require.NoError(t, s.DeleteVolumes(mc, s))

	mc.AssertExpectations(t)
}
