package service

import (
	"sync"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newContainerAndDockerTest(t *testing.T) (container.Container, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.Nil(t, err)

	return container, dt
}

func newFromServiceAndDockerTest(t *testing.T, s *Service) (*Service, *dockertest.Testing) {
	c, dt := newContainerAndDockerTest(t)
	s, err := FromService(s, ContainerOption(c))
	require.NoError(t, err)
	return s, dt
}

func newFromServiceAndContainerMocks(t *testing.T, s *Service) (*Service, *mocks.Container) {
	m := &mocks.Container{}
	s, err := FromService(s, ContainerOption(m))
	require.NoError(t, err)
	return s, m
}

func TestGenerateId(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestGenerateId",
	})
	hash := service.computeHash()
	require.Equal(t, "bb2239f3d1f52c4dffe268cbca5a43005b6c993a", hash)
}

func TestNoCollision(t *testing.T) {
	service1, _ := FromService(&Service{
		Name: "TestNoCollision",
	})

	service2, _ := FromService(&Service{
		Name: "TestNoCollision2",
	})

	require.NotEqual(t, service1.ID, service2.ID)
}

func TestNew(t *testing.T) {
	var (
		path = "../service-test/task"
		hash = "1"
	)

	mc := &mocks.Container{}
	mc.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := xarchive.GzippedTar(path)
	require.NoError(t, err)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		require.Equal(t, DeployStatus{
			Message: "Receiving service context...",
			Type:    DRunning,
		}, <-statuses)

		require.Equal(t, DeployStatus{
			Message: "Service context received with success.",
			Type:    DDonePositive,
		}, <-statuses)

		require.Equal(t, DeployStatus{
			Message: "Building Docker image...",
			Type:    DRunning,
		}, <-statuses)

		require.Equal(t, DeployStatus{
			Message: "Image built with success.",
			Type:    DDonePositive,
		}, <-statuses)
	}()

	s, err := New(archive,
		ContainerOption(mc),
		DeployStatusOption(statuses),
	)
	require.NoError(t, err)
	require.Equal(t, "service", s.Dependencies[0].Key)
	require.Equal(t, hash, s.Dependencies[0].Image)

	mc.AssertExpectations(t)
	wg.Wait()
}

func TestInjectDefinitionWithConfig(t *testing.T) {
	command := "xxx"
	s := &Service{}
	s.injectDefinition(&importer.ServiceDefinition{
		Configuration: &importer.Dependency{
			Command: command,
		},
	})
	require.Equal(t, command, s.configuration.Command)
}

func TestInjectDefinitionWithDependency(t *testing.T) {
	var (
		s     = &Service{}
		image = "xxx"
	)
	s.injectDefinition(&importer.ServiceDefinition{
		Dependencies: map[string]*importer.Dependency{
			"test": {
				Image: image,
			},
		},
	})
	require.Equal(t, s.Dependencies[0].Image, image)
}
