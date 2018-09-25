package service

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/stretchr/testify/require"
)

func newContainerAndDockerTest(t *testing.T) (*container.Container, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.NoError(t, err)

	return container, dt
}

func newFromServiceAndDockerTest(t *testing.T, s *Service) (*Service, *dockertest.Testing) {
	c, dt := newContainerAndDockerTest(t)
	s, err := FromService(s, ContainerOption(c))
	require.NoError(t, err)
	return s, dt
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
		hash = "sha256:x"
	)

	container, dt := newContainerAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"stream":"%s"}`, hash))), nil)

	archive, err := xarchive.GzippedTar(path)
	require.NoError(t, err)

	s, err := New(archive,
		ContainerOption(container),
	)
	require.NoError(t, err)
	require.Equal(t, "service", s.Dependencies[0].Key)
	require.Equal(t, hash, s.Dependencies[0].Image)
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
