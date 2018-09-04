package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
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

func TestStartService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(service.Dependencies), len(dockerServices))
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartWith2Dependencies(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartWith2Dependencies",
		Dependencies: []*Dependency{
			{
				Key:   "testa",
				Image: "nginx:stable-alpine",
			},
			{
				Key:   "testb",
				Image: "alpine:latest",
			},
		},
	}, ContainerOption(defaultContainer))
	servicesID, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.Dependencies
	container1, err1 := defaultContainer.FindContainer(deps[0].namespace())
	container2, err2 := defaultContainer.FindContainer(deps[1].namespace())
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, "nginx:stable-alpine", container1.Config.Image)
	require.Equal(t, "alpine:latest", container2.Config.Image)
}

func TestStartAgainService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartAgainService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestPartiallyRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestPartiallyRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "testa",
				Image: "nginx:stable-alpine",
			},
			{
				Key:   "testb",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	service.Dependencies[0].Stop()
	status, _ := service.Status()
	require.Equal(t, PARTIAL, status)
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, _ = service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartDependency(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	defer defaultContainer.DeleteNetwork(service.namespace())
	dep := service.Dependencies[0]
	serviceID, err := dep.Start(networkID)
	defer dep.Stop()
	require.Nil(t, err)
	require.NotEqual(t, "", serviceID)
	status, _ := dep.Status()
	require.Equal(t, container.RUNNING, status)
}

func TestNetworkCreated(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestNetworkCreated",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	network, err := defaultContainer.FindNetwork(service.namespace())
	require.Nil(t, err)
	require.NotEqual(t, "", network.ID)
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestStartStopStart(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartStopStart",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	service.Stop()
	time.Sleep(10 * time.Second)
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 1)
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestServiceDependenciesListensFromSamePort(t *testing.T) {
	var (
		service, _ = FromService(&Service{
			Name: "TestServiceDependenciesListensFromSamePort",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "nginx:stable-alpine",
					Ports: []string{"80"},
				},
			},
		}, ContainerOption(defaultContainer))

		service1, _ = FromService(&Service{
			Name: "TestServiceDependenciesListensFromSamePort1",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "nginx:stable-alpine",
					Ports: []string{"80"},
				},
			},
		}, ContainerOption(defaultContainer))
	)
	_, err := service.Start()
	require.NoError(t, err)
	defer service.Stop()

	_, err = service1.Start()
	require.NotZero(t, err)
	require.Contains(t, err.Error(), `port '80' is already in use`)
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
