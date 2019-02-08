package service

import (
	"bytes"
	"io"
	"testing"

	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/stretchr/testify/assert"
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
		{
			r:    new(bytes.Buffer),
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			env:  map[string]string{},
		},
		{
			r:    bytes.NewBufferString("a"),
			hash: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
			env:  map[string]string{},
		},
		{
			r:    new(bytes.Buffer),
			hash: "380918b946a526640a40df5dced6516794f3d97bbd9e6bb553d037c4439f31c3",
			env:  map[string]string{"": ""},
		},
		{
			r:    new(bytes.Buffer),
			hash: "3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a",
			env:  map[string]string{"foo": "bar"},
		},
		{
			r:    bytes.NewBufferString("a"),
			hash: "91ae92a2599b5df0200f0416f6abb44940b4535f5e072e3583f7ea26482b1df0",
			env:  map[string]string{"foo": "bar"},
		},
		{
			r:    new(bytes.Buffer),
			hash: "3d011e09502a84552a0f8ae112d024cc2c115597e3a577d5f49007902c221dc5",
			env:  map[string]string{"hello": "world"},
		},
		{
			r:    new(bytes.Buffer),
			hash: "2ca363bd7b02b93f3a7e7dd5850635d37793a5fd65ae63742e218cfca29cead5",
			env:  map[string]string{"foo": "bar", "hello": "world"},
		},
	}
	s := &Service{}
	for _, test := range tests {
		hash, err := s.computeHash(test.r, test.env)
		assert.NoError(t, err)
		assert.Equal(t, test.hash, hash)
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
