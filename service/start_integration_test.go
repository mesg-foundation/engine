// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStartServiceIntegration(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStartService",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)
	dockerServices, err := service.Start(c)
	defer service.Stop(c)
	require.NoError(t, err)
	require.Equal(t, len(service.Dependencies), len(dockerServices))
	status, _ := service.Status(c)
	require.Equal(t, RUNNING, status)
}

func TestIntegrationStartWith2DependenciesIntegration(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
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
		}
		c = newIntegrationContainer(t)
	)
	servicesID, err := service.Start(c)
	defer service.Stop(c)
	require.NoError(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.Dependencies
	container1, err1 := c.FindContainer(deps[0].namespace(service.namespace()))
	container2, err2 := c.FindContainer(deps[1].namespace(service.namespace()))
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, "http-server:latest", container1.Config.Image)
	require.Equal(t, "sleep:latest", container2.Config.Image)
}

func TestIntegrationStartAgainService(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStartAgainService",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	service.Start(c)
	defer service.Stop(c)
	dockerServices, err := service.Start(c)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := service.Status(c)
	require.Equal(t, RUNNING, status)
}

// // TODO: Disable this test in order to have the CI working
// func TestIntegrationPartiallyRunningService(t *testing.T) {
// 	var (
// 		service = &Service{
// 			Name: "TestPartiallyRunningService",
// 			Dependencies: []*Dependency{
// 				{
// 					Key:   "testa",
// 					Image: "http-server",
// 				},
// 				{
// 					Key:   "testb",
// 					Image: "http-server",
// 				},
// 			},
// 		}
// 		c = newIntegrationContainer(t)
// 	)

// 	service.Start(c)
// 	defer service.Stop(c)
// 	service.Dependencies[0].Stop(c)
// 	status, _ := service.Status(c)
// 	require.Equal(t, PARTIAL, status)
// 	dockerServices, err := service.Start(c)
// 	require.NoError(t, err)
// 	require.Equal(t, len(dockerServices), len(service.Dependencies))
// 	status, _ = service.Status(c)
// 	require.Equal(t, RUNNING, status)
// }

func TestIntegrationStartDependency(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStartDependency",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	networkID, err := c.CreateNetwork(service.namespace())
	require.NoError(t, err)
	defer c.DeleteNetwork(service.namespace())
	dep := service.Dependencies[0]
	serviceID, err := dep.Start(c, service, networkID)
	defer service.Stop(c)
	require.NoError(t, err)
	require.NotEqual(t, "", serviceID)
	status, _ := dep.Status(c, service)
	require.Equal(t, container.RUNNING, status)
}

func TestIntegrationNetworkCreated(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestNetworkCreated",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	service.Start(c)
	defer service.Stop(c)
	network, err := c.FindNetwork(service.namespace())
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestIntegrationStartStopStart(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStartStopStart",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	service.Start(c)
	service.Stop(c)
	dockerServices, err := service.Start(c)
	defer service.Stop(c)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), 1)
	status, _ := service.Status(c)
	require.Equal(t, RUNNING, status)
}

func TestIntegrationServiceDependenciesListensFromSamePort(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestServiceDependenciesListensFromSamePort",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}
		service1 = &Service{
			Hash: "2",
			Name: "TestServiceDependenciesListensFromSamePort1",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	_, err := service.Start(c)
	require.NoError(t, err)
	defer service.Stop(c)

	_, err = service1.Start(c)
	require.NotZero(t, err)
	require.Contains(t, err.Error(), `port '80' is already in use`)
}

func TestStartWithSamePorts(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStartWithSamePorts",
			Dependencies: []*Dependency{
				{
					Key:   "1",
					Image: "nginx",
					Ports: []string{"80"},
				},
				{
					Key:   "2",
					Image: "nginx",
					Ports: []string{"80"},
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	_, err := service.Start(c)
	require.Error(t, err)
}
