package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newFromServiceAndContainerMocks(t *testing.T, s *Service) (*Service, *mocks.Container) {
	m := &mocks.Container{}
	s, err := FromService(s, ContainerOption(m))
	require.NoError(t, err)
	return s, m
}

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		hash string
		env  map[string]string
	}{
		{hash: "4ef21a2e92a0bee1842cae888a408df5796683f4", env: map[string]string{}},
		{hash: "4ef21a2e92a0bee1842cae888a408df5796683f4", env: map[string]string{"": ""}},
		{hash: "c57d05deeb30464b209430b346d353147e18b2dd", env: map[string]string{"foo": "bar"}},
		{hash: "328aa77ff69d765da1398f4ee93503455925da24", env: map[string]string{"hello": "world"}},
		{hash: "c4cf57e9173296292ffc3498cac245e471b08e36", env: map[string]string{"foo": "bar", "hello": "world"}},
	}
	for _, test := range tests {
		service, _ := FromService(&Service{
			Name: "TestGenerateHash",
		})
		f, err := os.Open(filepath.Join("..", "service-test", "test.tar.gz"))
		require.NoError(t, err)
		defer f.Close()
		hash, err := service.computeHash(f, test.env)
		require.NoError(t, err)
		require.Equal(t, test.hash, hash)
	}
}

func TestNew(t *testing.T) {
	var (
		path = "../service-test/task"
		hash = "1"
	)

	mc := &mocks.Container{}
	mc.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := xarchive.GzippedTar(path, nil)
	require.NoError(t, err)

	statuses := make(chan DeployStatus, 4)

	s, err := New(archive, nil,
		ContainerOption(mc),
		DeployStatusOption(statuses),
	)
	require.NoError(t, err)
	require.Equal(t, "service", s.Dependencies[0].Key)
	require.Equal(t, hash, s.Dependencies[0].Image)
	require.Len(t, s.Dependencies[0].Env, 0)

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    DRunning,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success",
		Type:    DDonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    DRunning,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Image built with success",
		Type:    DDonePositive,
	}, <-statuses)

	mc.AssertExpectations(t)
}

func TestNewWithDefaultEnv(t *testing.T) {
	var (
		path = "../service-test/env"
		hash = "1"
		env  = []string{"A=1", "B=2"}
	)

	mc := &mocks.Container{}
	mc.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := xarchive.GzippedTar(path, nil)
	require.NoError(t, err)

	s, err := New(archive, nil,
		ContainerOption(mc),
	)
	require.NoError(t, err)
	require.Equal(t, "service", s.Dependencies[0].Key)
	require.Equal(t, hash, s.Dependencies[0].Image)
	require.Equal(t, env, s.Dependencies[0].Env)

	mc.AssertExpectations(t)
}

func TestNewWithOverwrittenEnv(t *testing.T) {
	var (
		path = "../service-test/env"
		hash = "1"
		env  = []string{"A=3", "B=4"}
	)

	mc := &mocks.Container{}
	mc.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := xarchive.GzippedTar(path, nil)
	require.NoError(t, err)

	s, err := New(archive, xos.EnvSliceToMap(env),
		ContainerOption(mc),
	)
	require.NoError(t, err)
	require.Equal(t, "service", s.Dependencies[0].Key)
	require.Equal(t, hash, s.Dependencies[0].Image)
	require.Equal(t, env, s.Dependencies[0].Env)

	mc.AssertExpectations(t)
}

func TestNewWitNotDefinedEnv(t *testing.T) {
	var (
		path = "../service-test/task"
	)

	mc := &mocks.Container{}

	archive, err := xarchive.GzippedTar(path, nil)
	require.NoError(t, err)

	_, err = New(archive, xos.EnvSliceToMap([]string{"A=1", "B=2"}),
		ContainerOption(mc),
	)
	require.Equal(t, ErrNotDefinedEnv{[]string{"A", "B"}}, err)

	mc.AssertExpectations(t)
}

func TestErrNotDefinedEnv(t *testing.T) {
	require.Equal(t, ErrNotDefinedEnv{[]string{"A", "B"}}.Error(),
		`environment variable(s) "A, B" not defined in mesg.yml (under configuration.env key)`)
}

func TestInjectDefinitionWithConfig(t *testing.T) {
	command := "xxx"
	s := &Service{}
	s.injectDefinition(&importer.ServiceDefinition{
		Configuration: &importer.Dependency{
			Command: command,
		},
	})
	require.Equal(t, command, s.configuration().Command)
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
