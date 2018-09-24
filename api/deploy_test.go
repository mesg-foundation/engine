package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/stretchr/testify/require"
)

func TestDeployService(t *testing.T) {
	path := filepath.Join("..", "service-test", "task")

	a, dt := newAPIAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := xarchive.GzippedTar(path)
		require.Nil(t, err)

		service, validationError, err := a.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.Nil(t, err)
		require.Len(t, service.ID, 40)
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success.",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "[DEPRECATED] Please use .dockerignore instead of .mesgignore",
		Type:    DoneNegative,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Image built with success.",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service deployed.",
		Type:    DonePositive,
	}, <-statuses)

	wg.Wait()
}

func TestDeployInvalidService(t *testing.T) {
	path := filepath.Join("..", "service-test", "invalid")

	a, dt := newAPIAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := xarchive.GzippedTar(path)
		require.Nil(t, err)

		service, validationError, err := a.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, service)
		require.Nil(t, err)
		require.Equal(t, (&importer.ValidationError{}).Error(), validationError.Error())
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success.",
		Type:    DonePositive,
	}, <-statuses)

	select {
	case <-statuses:
		t.Error("should not send further status messages")
	default:
	}

	wg.Wait()
}

func TestDeployServiceFromURL(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	a, dt := newAPIAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		service, validationError, err := a.DeployServiceFromURL(url, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.Nil(t, err)
		require.Len(t, service.ID, 40)
	}()

	require.Equal(t, DeployStatus{
		Message: "Downloading service...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service downloaded with success.",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success.",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Image built with success.",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service deployed.",
		Type:    DonePositive,
	}, <-statuses)

	wg.Wait()
}

func TestCreateTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, err := deployer.createTempDir()
	defer os.RemoveAll(path)
	require.Nil(t, err)
	require.NotEqual(t, "", path)
}

func TestRemoveTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, _ := deployer.createTempDir()
	err := os.RemoveAll(path)
	require.Nil(t, err)
}
