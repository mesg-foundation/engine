// +build integration

package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeleteVolumes(t *testing.T) {
	var (
		dependencyKey1 = "1"
		dependencyKey2 = "2"
		volumeA        = "/a"
		volumeB        = "/b"
		s, _           = FromService(&Service{
			Name: "TestIntegrationDeleteVolumes",
			Dependencies: []*Dependency{
				{
					Key:     dependencyKey1,
					Image:   "http-server",
					Volumes: []string{volumeA, volumeB},
				},
				{
					Key:         dependencyKey2,
					Image:       "http-server",
					VolumesFrom: []string{dependencyKey1},
				},
			},
		}, ContainerOption(newIntegrationContainer(t, container.TimeoutOption(time.Minute))))
	)
	_, err := s.Start()
	require.NoError(t, err)
	err = s.Stop()
	require.NoError(t, err)
	err = s.DeleteVolumes()
	require.NoError(t, err)
}
