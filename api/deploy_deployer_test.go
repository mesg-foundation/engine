package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, err := deployer.createTempDir()
	defer os.RemoveAll(path)
	require.Nil(t, err)
	require.NotEqual(t, "", path)
}

func TestRemoveTempFolder(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	deployer := newServiceDeployer(a)

	path, _ := deployer.createTempDir()
	err := os.RemoveAll(path)
	require.Nil(t, err)
}
