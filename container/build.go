package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

// Build a docker image
func Build(path string, namespace []string) (tag string, err error) {
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
	tag = ServiceTag(namespace)
	_, err = client.ImageBuild(context.Background(), buildContext, types.ImageBuildOptions{
		Tags:           []string{tag},
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	})
	return
}
