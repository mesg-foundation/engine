package container

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

// BuildResponse is the object that is returned by the docker api in json
type BuildResponse struct {
	Stream string `json:"stream"`
	Error  string `json:"error"`
}

// Build a docker image
func Build(path string) (tag string, err error) {
	buildContext, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return
	}
	defer buildContext.Close()
	client, err := Client()
	if err != nil {
		return
	}
	response, err := client.ImageBuild(context.Background(), buildContext, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	})
	if err != nil {
		return
	}
	tag, err = parseBuildResponse(response)
	return
}

func parseBuildResponse(response types.ImageBuildResponse) (tag string, err error) {
	lastOutput, err := extractLastOutputFromBuildResponse(response)
	if err != nil {
		return
	}
	var buildResponse BuildResponse
	err = json.Unmarshal([]byte(lastOutput), &buildResponse)
	if err != nil {
		err = errors.New("Could not parse container build response. " + err.Error())
		return
	}
	if buildResponse.Error != "" {
		err = errors.New("Image build failed. " + buildResponse.Error)
		return
	}
	tag = strings.TrimSuffix(buildResponse.Stream, "\n")
	return
}

func extractLastOutputFromBuildResponse(response types.ImageBuildResponse) (lastOutput string, err error) {
	defer response.Body.Close()
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	rs := strings.Split(string(r), "\n")
	i := len(rs) - 1
	for lastOutput == "" && i >= 0 {
		lastOutput = rs[i]
		i--
	}
	if lastOutput == "" {
		err = errors.New("Could not parse container build response")
		return
	}
	return
}
