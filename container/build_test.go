package container

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newImageBuildResponse(body string) types.ImageBuildResponse {
	return types.ImageBuildResponse{
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}
}

func TestBuild(t *testing.T) {
	var (
		expectedTag = "test:x"
		options     = types.ImageBuildOptions{
			Remove:         true,
			ForceRemove:    true,
			SuppressOutput: true,
			Tags:           []string{expectedTag},
		}
		c, m = newTesting(t)
	)

	imageBuildMatchContext := func(context io.Reader) bool {
		buildContext, err := archive.TarWithOptions("test/", &archive.TarOptions{
			Compression:     archive.Gzip,
			ExcludePatterns: []string{"ignoreme"},
		})
		if err != nil {
			return false
		}
		defer buildContext.Close()
		wantedData, err := ioutil.ReadAll(buildContext)
		if err != nil {
			return false
		}
		data, err := ioutil.ReadAll(context)
		if err != nil {
			return false
		}
		return bytes.Equal(wantedData, data)
	}

	var (
		imageBuildArguments = []interface{}{
			mock.Anything,
			mock.MatchedBy(imageBuildMatchContext),
			options,
		}
		imageBuildResponse = []interface{}{
			newImageBuildResponse(fmt.Sprintf(`{"stream":"%s\n"}`, expectedTag)),
			nil,
		}
		imageBuildRun = func(args mock.Arguments) {
			ioutil.ReadAll(args.Get(1).(io.Reader))
		}
	)

	m.On("ImageBuild", imageBuildArguments...).Once().Return(imageBuildResponse...).Run(imageBuildRun)

	tag, err := c.Build("test/", "test", "x")
	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)

	m.AssertExpectations(t)
}

func TestBuildResponseError(t *testing.T) {
	var (
		path                = "test-not-valid"
		c, m                = newTesting(t)
		imageBuildArguments = []interface{}{
			mock.Anything,
			mock.Anything,
			mock.Anything,
		}
		imageBuildResponse = []interface{}{
			newImageBuildResponse(`
{"stream":"Step 1/2 : FROM notExistingImage"}
{"stream":"\n"}
{"errorDetail":{"message": "invalid reference format: repository name must be lowercase"},"error":"invalid reference format: repository name must be lowercase"}`),
			nil,
		}
		imageBuildRun = func(args mock.Arguments) {
			ioutil.ReadAll(args.Get(1).(io.Reader))
		}
	)

	m.On("ImageBuild", imageBuildArguments...).Once().Return(imageBuildResponse...).Run(imageBuildRun)

	_, err := c.Build(path, "test", "x")
	require.Equal(t, "image build failed. invalid reference format: repository name must be lowercase", err.Error())

	m.AssertExpectations(t)
}

func TestBuildWrongPath(t *testing.T) {
	var (
		c, m                = newTesting(t)
		imageBuildArguments = []interface{}{
			mock.Anything,
			mock.Anything,
			mock.Anything,
		}
		imageBuildResponse = []interface{}{
			newImageBuildResponse(""),
			nil,
		}
		imageBuildRun = func(args mock.Arguments) {
			ioutil.ReadAll(args.Get(1).(io.Reader))
		}
	)

	m.On("ImageBuild", imageBuildArguments...).Once().Return(imageBuildResponse...).Run(imageBuildRun)

	_, err := c.Build("testss/", "test", "x")
	require.Equal(t, "could not parse container build response", err.Error())

	m.AssertExpectations(t)
}

func TestParseBuildResponse(t *testing.T) {
	err := parseBuildResponse(newImageBuildResponse("{\"stream\":\"ok\\n\"}"))
	require.NoError(t, err)
}
