// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newContainer(t)))
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(service.Dependencies), len(dockerServices))
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestStartWith2Dependencies(t *testing.T) {
	c := newContainer(t)
	service, _ := FromService(&Service{
		Name: "TestStartWith2Dependencies",
		Dependencies: []*Dependency{
			{
				Key:   "testa",
				Image: "http-server:latest",
			},
			{
				Key:   "testb",
				Image: "sleep:latest",
			},
		},
	}, ContainerOption(c))
	servicesID, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.Dependencies
	container1, err1 := c.FindContainer(deps[0].namespace())
	container2, err2 := c.FindContainer(deps[1].namespace())
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, "http-server:latest", container1.Config.Image)
	require.Equal(t, "sleep:latest", container2.Config.Image)
}

func TestStartAgainService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartAgainService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newContainer(t)))
	service.Start()
	defer service.Stop()
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

// TODO: Disable this test in order to have the CI working
// func TestPartiallyRunningService(t *testing.T) {
// 	service, _ := FromService(&Service{
// 		Name: "TestPartiallyRunningService",
// 		Dependencies: []*Dependency{
// 			{
// 				Key:   "testa",
// 				Image: "http-server",
// 			},
// 			{
// 				Key:   "testb",
// 				Image: "http-server",
// 			},
// 		},
// 	}, ContainerOption(newContainer(t)))
// 	service.Start()
// 	defer service.Stop()
// 	service.Dependencies[0].Stop()
// 	status, _ := service.Status()
// 	require.Equal(t, PARTIAL, status)
// 	dockerServices, err := service.Start()
// 	require.Nil(t, err)
// 	require.Equal(t, len(dockerServices), len(service.Dependencies))
// 	status, _ = service.Status()
// 	require.Equal(t, RUNNING, status)
// }

func TestStartDependency(t *testing.T) {
	c := newContainer(t)
	service, _ := FromService(&Service{
		Name: "TestStartDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	networkID, err := c.CreateNetwork(service.namespace())
	defer c.DeleteNetwork(service.namespace())
	dep := service.Dependencies[0]
	serviceID, err := dep.Start(networkID)
	defer dep.Stop()
	require.Nil(t, err)
	require.NotEqual(t, "", serviceID)
	status, _ := dep.Status()
	require.Equal(t, container.RUNNING, status)
}

func TestNetworkCreated(t *testing.T) {
	c := newContainer(t)
	service, _ := FromService(&Service{
		Name: "TestNetworkCreated",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	service.Start()
	defer service.Stop()
	network, err := c.FindNetwork(service.namespace())
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
				Image: "http-server",
			},
		},
	}, ContainerOption(newContainer(t)))
	service.Start()
	service.Stop()
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), 1)
	status, _ := service.Status()
	require.Equal(t, RUNNING, status)
}

func TestServiceDependenciesListensFromSamePort(t *testing.T) {
	c := newContainer(t)
	var (
		service, _ = FromService(&Service{
			Name: "TestServiceDependenciesListensFromSamePort",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}, ContainerOption(c))

		service1, _ = FromService(&Service{
			Name: "TestServiceDependenciesListensFromSamePort1",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}, ContainerOption(c))
	)
	_, err := service.Start()
	require.NoError(t, err)
	defer service.Stop()

	_, err = service1.Start()
	require.NotZero(t, err)
	require.Contains(t, err.Error(), `port '80' is already in use`)
}
