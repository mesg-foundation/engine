// +build integration

package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeleteVolumes(t *testing.T) {
	// TODO: the following test doesn't work on CircleCI because we don't use "machine"
	// https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-mount-volumes-to-docker-containers-
	t.Skip("doesn't work on CircleCI because we don't use machine")
	var (
		dependencyKey1 = "1"
		dependencyKey2 = "2"
		volumeA        = "/a"
		volumeB        = "/b"
		s              = &service.Service{
			Hash: hash.Int(1),
			Name: "TestIntegrationDeleteVolumes",
			Dependencies: []*service.Dependency{
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
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	_, err := m.Start(s)
	require.NoError(t, err)
	err = m.Stop(s)
	require.NoError(t, err)
	err = m.Delete(s)
	require.NoError(t, err)
}
