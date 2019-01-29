package service

import (
	"bytes"
	"io"
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
		r    io.Reader
		hash string
		env  map[string]string
	}{
		{r: new(bytes.Buffer), hash: "da39a3ee5e6b4b0d3255bfef95601890afd80709", env: map[string]string{}},
		{r: bytes.NewBufferString("a"), hash: "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8", env: map[string]string{}},
		{r: new(bytes.Buffer), hash: "21606782c65e44cac7afbb90977d8b6f82140e76", env: map[string]string{"": ""}},
		{r: new(bytes.Buffer), hash: "2fb8f40115dd1e695cbe23d4f97ce5b1fb697eee", env: map[string]string{"foo": "bar"}},
		{r: bytes.NewBufferString("a"), hash: "51420feb07a534f887c2839559e0bce5212c6b15", env: map[string]string{"foo": "bar"}},
		{r: new(bytes.Buffer), hash: "2d60d4e129a4a54062d8f982c397f56d11d7a9b9", env: map[string]string{"hello": "world"}},
		{r: new(bytes.Buffer), hash: "069d960caf60fde30133b8150e562c770d503ea9", env: map[string]string{"foo": "bar", "hello": "world"}},
	}
	s := &Service{}
	for _, test := range tests {
		hash, err := s.computeHash(test.r, test.env)
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
