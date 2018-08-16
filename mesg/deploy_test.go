package mesg

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
	path := "./service-test"

	mesg, dt := newMESGAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.Nil(t, err)

		service, validationError, err := mesg.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.Nil(t, err)
		require.Equal(t, 1, structhash.Version(service.Id))
	}()

	require.Equal(t, "Sending service context to core daemon...", <-statuses)
	require.Equal(t, "Building Docker image...", <-statuses)
	require.Equal(t, fmt.Sprintf("%s Completed.", aurora.Green("✔")), <-statuses)

	wg.Wait()
}

func TestDeployInvalidService(t *testing.T) {
	path := "./service-test-invalid"

	mesg, dt := newMESGAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := archive.TarWithOptions(path, &archive.TarOptions{
			Compression: archive.Gzip,
		})
		require.Nil(t, err)

		service, validationError, err := mesg.DeployService(archive, DeployServiceStatusOption(statuses))
		require.Nil(t, service)
		require.Nil(t, err)
		require.Equal(t, (&importer.ValidationError{}).Error(), validationError.Error())
	}()

	require.Equal(t, "Sending service context to core daemon...", <-statuses)

	select {
	case <-statuses:
		t.Error("should not send further status messages")
	default:
	}

	wg.Wait()
}

func TestDeployServiceFromURL(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	mesg, dt := newMESGAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		service, validationError, err := mesg.DeployServiceFromURL(url, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.Nil(t, err)
		require.Equal(t, 1, structhash.Version(service.Id))
	}()

	require.Equal(t, "Downloading service...", <-statuses)
	require.Equal(t, "Building Docker image...", <-statuses)
	require.Equal(t, fmt.Sprintf("%s Completed.", aurora.Green("✔")), <-statuses)

	wg.Wait()
}
