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
	service := &Service{
		Name: "TestStartService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(service.Dependencies), len(dockerServices))
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartWith2Dependencies(t *testing.T) {
	service := &Service{
		Name: "TestStartWith2Dependencies",
		Dependencies: map[string]*Dependency{
			"testa": {
				Image: "http-server:latest",
			},
			"testb": {
				Image: "sleep:latest",
			},
		},
	}
	servicesID, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.DependenciesFromService()
	container1, err1 := defaultContainer.FindContainer(deps[0].namespace())
	container2, err2 := defaultContainer.FindContainer(deps[1].namespace())
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, "http-server:latest", container1.Config.Image)
	require.Equal(t, "sleep:latest", container2.Config.Image)
}

func TestStartAgainService(t *testing.T) {
	service := &Service{
		Name: "TestStartAgainService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	service.Start()
	defer service.Stop()
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestPartiallyRunningService(t *testing.T) {
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]*Dependency{
			"testa": {
				Image: "http-server",
			},
			"testb": {
				Image: "http-server",
			},
		},
	}
	service.Start()
	defer service.Stop()
	service.DependenciesFromService()[0].Stop()
	status, _ := service.Status()
	require.Equal(t, PARTIAL, status)
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, _ = service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartDependency(t *testing.T) {
	service := &Service{
		Name: "TestStartDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	defer defaultContainer.DeleteNetwork(service.namespace())
	dep := service.DependenciesFromService()[0]
	serviceID, err := dep.Start(networkID)
	defer dep.Stop()
	require.Nil(t, err)
	require.NotEqual(t, "", serviceID)
	status, _ := dep.Status()
	require.Equal(t, container.RUNNING, status)
}

func TestNetworkCreated(t *testing.T) {
	service := &Service{
		Name: "TestNetworkCreated",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	service.Start()
	defer service.Stop()
	network, err := defaultContainer.FindNetwork(service.namespace())
	require.Nil(t, err)
	require.NotEqual(t, "", network.ID)
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestStartStopStart(t *testing.T) {
	service := &Service{
		Name: "TestStartStopStart",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
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
		service = &Service{
			Name: "TestServiceDependenciesListensFromSamePort",
			Dependencies: map[string]*Dependency{
				"test": {
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}

		service1 = &Service{
			Name: "TestServiceDependenciesListensFromSamePort1",
			Dependencies: map[string]*Dependency{
				"test": {
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}
	)
	_, err := service.Start()
	require.NoError(t, err)
	defer service.Stop()

	_, err = service1.Start()
	require.NotZero(t, err)
	require.Contains(t, err.Error(), `port '80' is already in use`)
}

func TestExtractVolumes(t *testing.T) {
	dep := &DependencyFromService{}
	_, err := dep.extractVolumes()
	require.NotNil(t, err)

	dep = &DependencyFromService{
		Name:    "test",
		Service: &Service{},
		Dependency: &Dependency{
			Volumes: []string{"foo", "bar"},
		},
	}
	volumes, err := dep.extractVolumes()
	require.Nil(t, err)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(dep.Service, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(dep.Service, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)

	dep = &DependencyFromService{
		Service: &Service{},
		Dependency: &Dependency{
			VolumesFrom: []string{"test"},
		},
	}
	_, err = dep.extractVolumes()
	require.NotNil(t, err)

	dep = &DependencyFromService{
		Service: &Service{
			Dependencies: map[string]*Dependency{
				"test": &Dependency{
					Volumes: []string{"foo", "bar"},
				},
			},
		},
		Dependency: &Dependency{
			VolumesFrom: []string{"test"},
		},
	}
	volumes, err = dep.extractVolumes()
	require.Nil(t, err)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(dep.Service, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(dep.Service, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)
}
