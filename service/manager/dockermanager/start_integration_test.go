// +build integration

package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStartServiceIntegration(t *testing.T) {
	var (
		s = &service.Service{
			Hash: "1",
			Name: "TestStartService",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)
	dockerServices, err := m.Start(s)
	defer m.Stop(s)
	require.NoError(t, err)
	require.Equal(t, len(s.Dependencies), len(dockerServices))
	status, _ := m.Status(s)
	require.Equal(t, service.RUNNING, status)
}

func TestIntegrationStartWith2DependenciesIntegration(t *testing.T) {
	var (
		service = &service.Service{
			Hash: "1",
			Name: "TestStartWith2Dependencies",
			Dependencies: []*service.Dependency{
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
		m = New(c)
	)
	servicesID, err := m.Start(service)
	defer m.Stop(service)
	require.NoError(t, err)
	require.Equal(t, 2, len(servicesID))
	deps := service.Dependencies
	container1, err1 := c.FindContainer(dependencyNamespace(serviceNamespace(service.Hash), deps[0].Key))
	container2, err2 := c.FindContainer(dependencyNamespace(serviceNamespace(service.Hash), deps[1].Key))
	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Equal(t, "http-server:latest", container1.Config.Image)
	require.Equal(t, "sleep:latest", container2.Config.Image)
}

func TestIntegrationStartAgainService(t *testing.T) {
	var (
		s = &service.Service{
			Hash: "1",
			Name: "TestStartAgainService",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	m.Start(s)
	defer m.Stop(s)
	dockerServices, err := m.Start(s)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	status, _ := m.Status(s)
	require.Equal(t, service.RUNNING, status)
}

// // TODO: Disable this test in order to have the CI working
// func TestIntegrationPartiallyRunningService(t *testing.T) {
// 	var (
// 		service = &Service{
// 			Name: "TestPartiallyRunningService",
// 			Dependencies: []*service.Dependency{
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
//		m = New(c)
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

func TestIntegrationNetworkCreated(t *testing.T) {
	var (
		service = &service.Service{
			Hash: "1",
			Name: "TestNetworkCreated",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	m.Start(service)
	defer m.Stop(service)
	network, err := c.FindNetwork(serviceNamespace(service.Hash))
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestIntegrationStartStopStart(t *testing.T) {
	var (
		s = &service.Service{
			Hash: "1",
			Name: "TestStartStopStart",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	m.Start(s)
	m.Stop(s)
	dockerServices, err := m.Start(s)
	defer m.Stop(s)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), 1)
	status, _ := m.Status(s)
	require.Equal(t, service.RUNNING, status)
}

func TestIntegrationServiceDependenciesListensFromSamePort(t *testing.T) {
	var (
		s = &service.Service{
			Hash: "1",
			Name: "TestServiceDependenciesListensFromSamePort",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}
		s1 = &service.Service{
			Hash: "2",
			Name: "TestServiceDependenciesListensFromSamePort1",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
					Ports: []string{"80"},
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	_, err := m.Start(s)
	require.NoError(t, err)
	defer m.Stop(s)

	_, err = m.Start(s1)
	require.NotZero(t, err)
	require.Contains(t, err.Error(), `port '80' is already in use`)
}

func TestStartWithSamePorts(t *testing.T) {
	var (
		service = &service.Service{
			Hash: "1",
			Name: "TestStartWithSamePorts",
			Dependencies: []*service.Dependency{
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
		m = New(c)
	)

	_, err := m.Start(service)
	require.Error(t, err)
}
