package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDockerFile(t *testing.T) {
	data, err := readDockerfile("./tests/service-valid")
	require.Nil(t, err)
	require.True(t, (len(data) > 0))
}

func TestReadDockerFileDoesNotExist(t *testing.T) {
	data, err := readDockerfile("./tests/docker-missing")
	require.NotNil(t, err)
	require.True(t, (len(data) == 0))
}
