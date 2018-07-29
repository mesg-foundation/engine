package container

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stvp/assert"
)

func TestBuild(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	tag, err := c.Build("test/")
	assert.Nil(t, err)
	assert.NotEqual(t, "", tag)
}

func TestBuildNotWorking(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	tag, err := c.Build("test-not-valid/")
	assert.NotNil(t, err)
	assert.Equal(t, "", tag)
}

func TestBuildWrongPath(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.Build("testss/")
	assert.NotNil(t, err)
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
