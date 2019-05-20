package api

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeployService(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "task")
		hash  = "1"
		a, at = newTesting(t)
	)
	defer at.close()

	at.containerMock.On("Build", mock.Anything).Once().Return(hash, nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.NoError(t, err)

		s, validationError, err := a.DeployService(archive, nil, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.NoError(t, err)
		require.Equal(t, "service", s.Configuration.Key)
		require.Equal(t, hash, s.Configuration.Image)
		require.Len(t, s.Configuration.Env, 0)
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Image built with success",
		Type:    DonePositive,
	}, <-statuses)

	wg.Wait()
	at.containerMock.AssertExpectations(t)
}

func TestDeployWithDefaultEnv(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "env")
		hash  = "1"
		env   = []string{"A=1", "B=2"}
		a, at = newTesting(t)
	)
	defer at.close()

	at.containerMock.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)

	s, validationError, err := a.DeployService(archive, nil)
	require.Nil(t, validationError)
	require.NoError(t, err)
	require.Equal(t, "service", s.Configuration.Key)
	require.Equal(t, hash, s.Configuration.Image)
	require.Equal(t, env, s.Configuration.Env)

	at.containerMock.AssertExpectations(t)
}

func TestDeployWithOverwrittenEnv(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "env")
		hash  = "1"
		env   = []string{"A=3", "B=4"}
		a, at = newTesting(t)
	)
	defer at.close()

	at.containerMock.On("Build", mock.Anything).Once().Return(hash, nil)

	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)

	s, validationError, err := a.DeployService(archive, xos.EnvSliceToMap(env))
	require.Nil(t, validationError)
	require.NoError(t, err)
	require.Equal(t, "service", s.Configuration.Key)
	require.Equal(t, hash, s.Configuration.Image)
	require.Equal(t, env, s.Configuration.Env)

	at.containerMock.AssertExpectations(t)
}

func TestDeployWitNotDefinedEnv(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "task")
		a, at = newTesting(t)
	)
	defer at.close()

	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)

	_, validationError, err := a.DeployService(archive, xos.EnvSliceToMap([]string{"A=1", "B=2"}))
	require.Nil(t, validationError)
	require.Equal(t, service.ErrNotDefinedEnv{[]string{"A", "B"}}, err)

	at.containerMock.AssertExpectations(t)
}

func TestDeployInvalidService(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "invalid")
		a, at = newTesting(t)
	)
	defer at.close()

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.NoError(t, err)

		_, validationError, err := a.DeployService(archive, nil, DeployServiceStatusOption(statuses))
		require.NoError(t, err)
		require.Equal(t, (&importer.ValidationError{}).Error(), validationError.Error())
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	wg.Wait()
	at.containerMock.AssertExpectations(t)
}

func TestDeployServiceFromURL(t *testing.T) {
	var (
		url   = "git://github.com/mesg-foundation/service-webhook#single-outputs"
		a, at = newTesting(t)
	)
	defer at.close()

	at.containerMock.On("Build", mock.Anything).Once().Return("1", nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, validationError, err := a.DeployServiceFromURL(url, nil, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.NoError(t, err)
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Image built with success",
		Type:    DonePositive,
	}, <-statuses)

	wg.Wait()
	at.containerMock.AssertExpectations(t)
}
