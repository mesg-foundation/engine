package serialize

import (
	"testing"

	"github.com/stvp/assert"
)

func TestReadDockerFile(t *testing.T) {
	data, err := readDockerfile("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, (len(data) > 0))
}

func TestReadDockerFileDoesNotExist(t *testing.T) {
	data, err := readDockerfile("./tests/docker-missing")
	assert.NotNil(t, err)
	assert.True(t, (len(data) == 0))
}
