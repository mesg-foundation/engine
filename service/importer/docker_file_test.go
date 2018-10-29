package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDockerFile(t *testing.T) {
	data, err := readDockerfile("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, (len(data) > 0))
}

func TestReadDockerFileDoesNotExist(t *testing.T) {
	data, err := readDockerfile("./tests/docker-missing")
	require.Error(t, err)
	require.True(t, (len(data) == 0))
}
