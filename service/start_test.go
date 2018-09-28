package service

import (
	"strconv"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/x/xnet"
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
		networkID          = "3"
		sharedNetworkID    = "4"
		s, mc              = newFromServiceAndContainerMocks(t, &Service{
			Name: serviceName,
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		})
	)

	var (
		d, _       = s.getDependency(dependencyKey)
		c, _       = config.Global()
		_, port, _ = xnet.SplitHostPort(c.Server.Address)
		endpoint   = c.Core.Name + ":" + strconv.Itoa(port)
		mounts, _  = d.extractVolumes()
	)

	mc.On("Status", d.namespace()).Once().Return(container.STOPPED, nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Once().Return(sharedNetworkID, nil)
	mc.On("StartService", container.ServiceOptions{
		Namespace: d.namespace(),
		Labels: map[string]string{
			"mesg.service": d.service.Name,
			"mesg.hash":    d.service.ID,
			"mesg.core":    c.Core.Name,
		},
		Image: d.Image,
		Args:  strings.Fields(d.Command),
		Env: container.MapToEnv(map[string]string{
			"MESG_TOKEN":        d.service.ID,
			"MESG_ENDPOINT":     endpoint,
			"MESG_ENDPOINT_TCP": endpoint,
		}),
		Mounts:     mounts,
		Ports:      d.extractPorts(),
		NetworksID: []string{networkID, sharedNetworkID},
	}).Once().Return(containerServiceID, nil)

	serviceIDs, err := s.Start()
	require.NoError(t, err)
	require.Len(t, serviceIDs, 1)
	require.Equal(t, containerServiceID, serviceIDs[0])

	mc.AssertExpectations(t)
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

func TestPartiallyRunningService(t *testing.T) {
	var (
		dependencyKey       = "1"
		dependencyKey2      = "2"
		networkID           = "3"
		sharedNetworkID     = "4"
		containerServiceIDs = []string{"5", "6"}
		s, mc               = newFromServiceAndContainerMocks(t, &Service{
			Name: "TestPartiallyRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
				{
					Key:   dependencyKey2,
					Image: "http-server",
				},
			},
		})
	)

	var (
		d, _       = s.getDependency(dependencyKey)
		d2, _      = s.getDependency(dependencyKey2)
		c, _       = config.Global()
		_, port, _ = xnet.SplitHostPort(c.Server.Address)
		endpoint   = c.Core.Name + ":" + strconv.Itoa(port)
		mounts, _  = d.extractVolumes()
	)

	mc.On("Status", d.namespace()).Return(container.STOPPED, nil)
	mc.On("Status", d2.namespace()).Return(container.RUNNING, nil)
	mc.On("StopService", d2.namespace()).Once().Return(nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Twice().Return(sharedNetworkID, nil)

	for i, d := range []*Dependency{d, d2} {
		mc.On("StartService", container.ServiceOptions{
			Namespace: d.namespace(),
			Labels: map[string]string{
				"mesg.service": d.service.Name,
				"mesg.hash":    d.service.ID,
				"mesg.core":    c.Core.Name,
			},
			Image: d.Image,
			Args:  strings.Fields(d.Command),
			Env: container.MapToEnv(map[string]string{
				"MESG_TOKEN":        d.service.ID,
				"MESG_ENDPOINT":     endpoint,
				"MESG_ENDPOINT_TCP": endpoint,
			}),
			Mounts:     mounts,
			Ports:      d.extractPorts(),
			NetworksID: []string{networkID, sharedNetworkID},
		}).Once().Return(containerServiceIDs[i], nil)
	}

	serviceIDs, err := s.Start()
	require.NoError(t, err)
	require.Len(t, serviceIDs, len(s.Dependencies))

	for i := range s.Dependencies {
		require.Equal(t, containerServiceIDs[i], serviceIDs[i])
	}

	mc.AssertExpectations(t)
}
