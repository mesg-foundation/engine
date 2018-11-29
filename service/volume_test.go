package service

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/require"
)

func TestCreateVolumes(t *testing.T) {
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
		d1, _      = s.getDependency(dependencyKey1)
		d2, _      = s.getDependency(dependencyKey2)
		mounts1, _ = d1.extractVolumes()
		mounts2, _ = d2.extractVolumes()
	)

	mc.On("CreateVolume", mounts1[0].Source).Once().Return(types.Volume{}, nil)
	mc.On("CreateVolume", mounts1[1].Source).Once().Return(types.Volume{}, nil)
	mc.On("CreateVolume", mounts2[0].Source).Once().Return(types.Volume{}, nil)
	mc.On("CreateVolume", mounts2[1].Source).Once().Return(types.Volume{}, nil)

	require.NoError(t, s.CreateVolumes())

	mc.AssertExpectations(t)
}

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
		d1, _      = s.getDependency(dependencyKey1)
		d2, _      = s.getDependency(dependencyKey2)
		mounts1, _ = d1.extractVolumes()
		mounts2, _ = d2.extractVolumes()
	)

	mc.On("DeleteVolume", mounts1[0].Source).Once().Return(nil)
	mc.On("DeleteVolume", mounts1[1].Source).Once().Return(nil)
	mc.On("DeleteVolume", mounts2[0].Source).Once().Return(nil)
	mc.On("DeleteVolume", mounts2[1].Source).Once().Return(nil)

	require.NoError(t, s.DeleteVolumes())

	mc.AssertExpectations(t)
}
