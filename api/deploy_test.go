package api

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/stretchr/testify/require"
)

func TestDeployService(t *testing.T) {
	path := "../service-test/task"

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
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Service context received with success.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s [DEPRECATED] Please use .dockerignore instead of .mesgignore", aurora.Red("⨯")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Image built with success.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	wg.Wait()
}

func TestDeployInvalidService(t *testing.T) {
	path := "../service-test/invalid"

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
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Service context received with success.", aurora.Green("✔")),
		Type:    DONE,
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
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Service downloaded with success.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Service context received with success.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Building Docker image...",
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Image built with success.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    DONE,
	}, <-statuses)

	wg.Wait()
}
