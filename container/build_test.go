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
		expectedTag = "sha256:1f6359c933421f53a7ef9e417bfa51b1c313c54878fdeb16de827f427e16d836"
		options     = types.ImageBuildOptions{
			Remove:         true,
			ForceRemove:    true,
			SuppressOutput: true,
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

	tag, err := c.Build("test/")
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

	tag, err := c.Build(path)
	require.Empty(t, tag)
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

	_, err := c.Build("testss/")
	require.Equal(t, "could not parse container build response", err.Error())

	m.AssertExpectations(t)
}

func TestParseBuildResponse(t *testing.T) {
	tag, err := parseBuildResponse(newImageBuildResponse("{\"stream\":\"ok\\n\"}"))
	require.NoError(t, err)
	require.Equal(t, tag, "ok")
}
