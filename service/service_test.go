package service

import (
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
