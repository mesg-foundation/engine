package service

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/stretchr/testify/require"
)

func TestExtractPortEmpty(t *testing.T) {
	dep := Dependency{}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	dep := &Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 2)
	require.Equal(t, ports[0].Target, uint32(80))
	require.Equal(t, ports[0].Published, uint32(80))
	require.Equal(t, ports[1].Target, uint32(8080))
	require.Equal(t, ports[1].Published, uint32(3000))
}

func TestExtractVolumes(t *testing.T) {
	s, _ := FromService(&Service{
		Dependencies: []*Dependency{{
			Key:     "test",
			Volumes: []string{"foo", "bar"},
		}},
	})
	volumes, err := s.Dependencies[0].extractVolumes()
	require.Nil(t, err)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(s, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(s, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)

	s, _ = FromService(&Service{
		Dependencies: []*Dependency{{
			VolumesFrom: []string{"test"},
		}},
	})
	_, err = s.Dependencies[0].extractVolumes()
	require.Error(t, err)

	s, _ = FromService(&Service{
		Dependencies: []*Dependency{
			{
				Key:     "test",
				Volumes: []string{"foo", "bar"},
			},
			{
				VolumesFrom: []string{"test"},
			}},
	})
	volumes, err = s.Dependencies[1].extractVolumes()
	require.Nil(t, err)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(s, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(s, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)
}

func TestStartService(t *testing.T) {
	var (
		containerServiceID = "1"
		dependencyKey      = "2"
		serviceName        = "TestStartService"
		s, dt              = newFromServiceAndDockerTest(t, &Service{
			Name: serviceName,
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)
	dt.ProvideContainerList(nil, dockertest.NotFoundErr{})
	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})
	dt.ProvideNetworkInspect(types.NetworkResource{ID: "3"}, nil)
	dt.ProvideNetworkInspect(types.NetworkResource{ID: "4"}, nil)
	// service create.
	dt.ProvideServiceCreate(types.ServiceCreateResponse{ID: containerServiceID}, nil)
	dockerServices, err := s.Start()
	require.NoError(t, err)
	require.Len(t, dockerServices, 1)
	require.Equal(t, containerServiceID, dockerServices[0])
	lc := <-dt.LastServiceCreate()
	require.Equal(t, types.ServiceCreateOptions{}, lc.Options)
	require.Equal(t, s.docker.Namespace([]string{s.ID, dependencyKey}), lc.Service.Name)
}
func TestStartWith2Dependencies(t *testing.T) {
	var (
		containerServiceID  = "1"
		containerServiceID2 = "2"
		dependencyKey       = "3"
		dependencyKey2      = "4"
		dependencyImage     = "5"
		dependencyImage2    = "6"
		serviceName         = "TestStartWith2Dependencies"
		s, dt               = newFromServiceAndDockerTest(t, &Service{
			Name: serviceName,
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: dependencyImage,
				},
				{
					Key:   dependencyKey2,
					Image: dependencyImage2,
				},
			},
		})
	)
	// for dep1 & dep2
	for i := 0; i < 2; i++ {
		dt.ProvideContainerList(nil, dockertest.NotFoundErr{})
		dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})
		dt.ProvideNetworkInspect(types.NetworkResource{ID: "3"}, nil)
		dt.ProvideNetworkInspect(types.NetworkResource{ID: "4"}, nil)
	}
	// service create.
	dt.ProvideServiceCreate(types.ServiceCreateResponse{ID: containerServiceID}, nil)
	dt.ProvideServiceCreate(types.ServiceCreateResponse{ID: containerServiceID2}, nil)
	servicesIDs, err := s.Start()
	require.NoError(t, err)
	require.Len(t, servicesIDs, 2)
	require.True(t, xstrings.SliceContains(servicesIDs, containerServiceID))
	require.True(t, xstrings.SliceContains(servicesIDs, containerServiceID2))
	images := []string{dependencyImage, dependencyImage2}
	for i := 0; i < 2; i++ {
		lc := <-dt.LastServiceCreate()
		require.True(t, xstrings.SliceContains(images, lc.Service.TaskTemplate.ContainerSpec.Image))
	}
}

func TestStartServiceRunning(t *testing.T) {
	var (
		s, dt = newFromServiceAndDockerTest(t, &Service{
			Dependencies: []*Dependency{
				{
					Key:   "1",
					Image: "2",
				},
			},
		})
	)

	dt.ProvideContainerList([]types.Container{{ID: "1"}}, nil)
	dt.ProvideContainerInspect(types.ContainerJSON{ContainerJSONBase: &types.ContainerJSONBase{ID: "1"}}, nil)
	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, nil)

	dockerServices, err := s.Start()
	require.NoError(t, err)
	require.Len(t, dockerServices, 0)
}
