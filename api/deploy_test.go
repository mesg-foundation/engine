// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	a, dt, closer := newAPIAndDockerTest(t)
	defer closer()
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := xarchive.GzippedTar(path, nil)
		require.NoError(t, err)

		service, validationError, err := a.DeployService(archive, nil, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.NoError(t, err)
		require.Len(t, service.Hash, 40)
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success",
		Type:    DonePositive,
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
}

func TestDeployInvalidService(t *testing.T) {
	path := filepath.Join("..", "service-test", "invalid")

	a, dt, closer := newAPIAndDockerTest(t)
	defer closer()
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		archive, err := xarchive.GzippedTar(path, nil)
		require.NoError(t, err)

		service, validationError, err := a.DeployService(archive, nil, DeployServiceStatusOption(statuses))
		require.Nil(t, service)
		require.NoError(t, err)
		require.Equal(t, (&importer.ValidationError{}).Error(), validationError.Error())
	}()

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success",
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

	a, dt, closer := newAPIAndDockerTest(t)
	defer closer()
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	statuses := make(chan DeployStatus)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		service, validationError, err := a.DeployServiceFromURL(url, nil, DeployServiceStatusOption(statuses))
		require.Nil(t, validationError)
		require.NoError(t, err)
		require.Len(t, service.Hash, 40)
	}()

	require.Equal(t, DeployStatus{
		Message: "Downloading service...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service downloaded with success",
		Type:    DonePositive,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Receiving service context...",
		Type:    Running,
	}, <-statuses)

	require.Equal(t, DeployStatus{
		Message: "Service context received with success",
		Type:    DonePositive,
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
}

func TestCreateTempFolder(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	deployer := newServiceDeployer(a, nil)

	path, err := deployer.createTempDir()
	defer os.RemoveAll(path)
	require.NoError(t, err)
	require.NotEqual(t, "", path)
}

func TestRemoveTempFolder(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	deployer := newServiceDeployer(a, nil)

	path, _ := deployer.createTempDir()
	err := os.RemoveAll(path)
	require.NoError(t, err)
}
