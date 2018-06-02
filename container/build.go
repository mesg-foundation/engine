package container

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

// Build a docker image
func Build(path string) (tag string, err error) {
	buildContext, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return
	}
	defer buildContext.Close()
	if err != nil {
		return
	}
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
	defer response.Body.Close()

	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	type BuildResponse struct {
		Stream string `json:"stream"`
	}
	var buildResponse BuildResponse
	err = json.Unmarshal(r, &buildResponse)
	if err != nil {
		return
	}
	tag = strings.TrimSuffix(buildResponse.Stream, "\n")
	return
}
