package api

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeployService(t *testing.T) {
	var (
		path  = filepath.Join("..", "service-test", "task")
		a, at = newTesting(t)
	)
	defer at.close()

	at.containerMock.On("Build", mock.Anything).Once().Return("1", nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := xarchive.GzippedTar(path, nil)
		require.NoError(t, err)

		_, validationError, err := a.DeployService(archive, nil, DeployServiceStatusOption(statuses))
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

		archive, err := xarchive.GzippedTar(path, nil)
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
		url   = "https://github.com/mesg-foundation/service-webhook.git"
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
