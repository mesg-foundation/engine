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

package container

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
)

// BuildResponse is the object that is returned by the docker api in json
type BuildResponse struct {
	Stream string `json:"stream"`
	Error  string `json:"error"`
}

// Build builds a docker image.
func (c *DockerContainer) Build(path string) (tag string, err error) {
	excludeFiles, err := dockerignoreFiles(path)
	if err != nil {
		return "", err
	}

	buildContext, err := xarchive.GzippedTar(path, excludeFiles)
	if err != nil {
		return "", err
	}
	defer buildContext.Close()
	response, err := c.client.ImageBuild(context.Background(), buildContext, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	})
	if err != nil {
		return "", err
	}
	return parseBuildResponse(response)
}

// dockerignoreFiles reads excluded files from .dockerignore.
func dockerignoreFiles(path string) ([]string, error) {
	f, err := os.Open(filepath.Join(path, ".dockerignore"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	return dockerignore.ReadAll(f)
}

func parseBuildResponse(response types.ImageBuildResponse) (tag string, err error) {
	lastOutput, err := extractLastOutputFromBuildResponse(response)
	if err != nil {
		return "", err
	}
	var buildResponse BuildResponse

	if err := json.Unmarshal([]byte(lastOutput), &buildResponse); err != nil {
		return "", fmt.Errorf("could not parse container build response. %s", err)
	}
	if buildResponse.Error != "" {
		return "", fmt.Errorf("image build failed. %s", buildResponse.Error)
	}
	return strings.TrimSuffix(buildResponse.Stream, "\n"), nil
}

func extractLastOutputFromBuildResponse(response types.ImageBuildResponse) (lastOutput string, err error) {
	defer response.Body.Close()
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	lastOutput = ""
	rs := strings.Split(string(r), "\n")
	i := len(rs) - 1
	for lastOutput == "" && i >= 0 {
		lastOutput = rs[i]
		i--
	}
	if lastOutput == "" {
		return "", errors.New("could not parse container build response")
	}
	return lastOutput, nil
}
