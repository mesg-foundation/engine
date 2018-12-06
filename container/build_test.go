package container

import (
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

func TestBuild(t *testing.T) {
	var (
		tag   = "sha256:1f6359c933421f53a7ef9e417bfa51b1c313c54878fdeb16de827f427e16d836"
		c, mc = newTesting(t)
	)

	mc.On("ImageBuild", mock.Anything, mock.Anything, types.ImageBuildOptions{
		Remove:         true,
		ForceRemove:    true,
		SuppressOutput: true,
	}).Once().
		Return(types.ImageBuildResponse{
			Body: ioutil.NopCloser(strings.NewReader(
				fmt.Sprintf(`{"stream":"%s\n"}`, tag),
			)),
		}, nil).
		Run(func(args mock.Arguments) {
			buildContext, err := archive.TarWithOptions("test/", &archive.TarOptions{
				Compression:     archive.Gzip,
				ExcludePatterns: []string{"ignoreme"},
			})
			require.NoError(t, err)
			defer buildContext.Close()

			wantedData, err := ioutil.ReadAll(buildContext)
			require.NoError(t, err)

			data, err := ioutil.ReadAll(args.Get(1).(io.Reader))
			require.NoError(t, err)
			require.Equal(t, wantedData, data)
		})

	tag1, err := c.Build("test/")
	require.NoError(t, err)
	require.Equal(t, tag, tag1)

	mc.AssertExpectations(t)
}

func TestBuildNotWorking(t *testing.T) {
	var (
		path  = "test-not-valid"
		c, mc = newTesting(t)
	)

	mc.On("ImageBuild", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(types.ImageBuildResponse{
			Body: ioutil.NopCloser(strings.NewReader(`
		{"stream":"Step 1/2 : FROM notExistingImage"}
		{"stream":"\n"}
		{"errorDetail":{"message":"invalid reference format: repository name must be lowercase"},"error":"invalid reference format: repository name must be lowercase"}`)),
		}, nil).
		Run(func(args mock.Arguments) {
			_, err := ioutil.ReadAll(args.Get(1).(io.Reader))
			require.NoError(t, err)
		})

	tag, err := c.Build(path)
	require.Equal(t, "image build failed. invalid reference format: repository name must be lowercase", err.Error())
	require.Equal(t, "", tag)

	mc.AssertExpectations(t)
}

func TestBuildWrongPath(t *testing.T) {
	var (
		c, mc = newTesting(t)
	)

	mc.On("ImageBuild", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(types.ImageBuildResponse{
			Body: ioutil.NopCloser(strings.NewReader("")),
		}, nil).
		Run(func(args mock.Arguments) {
			_, err := ioutil.ReadAll(args.Get(1).(io.Reader))
			require.NoError(t, err)
		})

	_, err := c.Build("testss/")
	require.Equal(t, "could not parse container build response", err.Error())

	mc.AssertExpectations(t)
}

func TestParseBuildResponseInvalidJSON(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("invalidJSON"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	_, err := parseBuildResponse(response)
	require.Error(t, err)
}

func TestParseBuildResponse(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("{\"stream\":\"ok\"}"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	tag, err := parseBuildResponse(response)
	require.NoError(t, err)
	require.Equal(t, tag, "ok")
}

func TestParseBuildResponseWithNewLine(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader("{\"stream\":\"ok\\n\"}"))
	response := types.ImageBuildResponse{
		Body: body,
	}
	tag, err := parseBuildResponse(response)
	require.NoError(t, err)
	require.Equal(t, tag, "ok")
}
