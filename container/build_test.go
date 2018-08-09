package container

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestBuild(t *testing.T) {
	path := "test/"
	tag := "sha256:1f6359c933421f53a7ef9e417bfa51b1c313c54878fdeb16de827f427e16d836"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(
		fmt.Sprintf(`{"stream":"%s\n"}`, tag),
	)), nil)

	tag1, err := c.Build(path)
	assert.Nil(t, err)
	assert.Equal(t, tag, tag1)

	li := <-dt.LastImageBuild()
	assert.True(t, len(li.FileData) > 0)
	assert.Equal(t, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	}, li.Options)
}

func TestBuildNotWorking(t *testing.T) {
	path := "test-not-valid/"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`
{"stream":"Step 1/2 : FROM notExistingImage"}
{"stream":"\n"}
{"errorDetail":{"message":"invalid reference format: repository name must be lowercase"},"error":"invalid reference format: repository name must be lowercase"}`)), nil)

	tag, err := c.Build(path)
	assert.Equal(t, "Image build failed. invalid reference format: repository name must be lowercase", err.Error())
	assert.Equal(t, "", tag)
}

func TestBuildWrongPath(t *testing.T) {
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader("")), nil)

	_, err := c.Build("testss/")
	assert.Equal(t, "Could not parse container build response", err.Error())
}

func TestParseBuildResponseInvalidJSON(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("invalidJSON"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	_, err := parseBuildResponse(response)
	assert.NotNil(t, err)
}

func TestParseBuildResponse(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("{\"stream\":\"ok\"}"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	tag, err := parseBuildResponse(response)
	assert.Nil(t, err)
	assert.Equal(t, tag, "ok")
}

func TestParseBuildResponseWithNewLine(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("{\"stream\":\"ok\\n\"}"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	tag, err := parseBuildResponse(response)
	assert.Nil(t, err)
	assert.Equal(t, tag, "ok")
}
