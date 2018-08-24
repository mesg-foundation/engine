package api

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"github.com/cnf/structhash"
	"github.com/docker/docker/pkg/archive"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service/importer"
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

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.Nil(t, err)

		service, validationError, err := a.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.Nil(t, err)
		require.Equal(t, 1, structhash.Version(service.Id))
	}()

	require.Equal(t, DeployStatus{
		Message: "Sending service context to core daemon...",
		Type:    RUNNING,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: fmt.Sprintf("%s Service context sent to core daemon with success.", aurora.Green("✔")),
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

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.Nil(t, err)

		service, validationError, err := a.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, service)
		require.Nil(t, err)
		require.Equal(t, (&importer.ValidationError{}).Error(), validationError.Error())
	}()

	require.Equal(t, "Sending service context to core daemon...", (<-statuses).Message)
	require.Equal(t, fmt.Sprintf("%s Service context sent to core daemon with success.", aurora.Green("✔")),
		(<-statuses).Message)

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
		require.Equal(t, 1, structhash.Version(service.Id))
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
