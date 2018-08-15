package config

import (
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestCreateConfigPath(t *testing.T) {
	err := createConfigPath()
	require.Nil(t, err)
}

func TestConfigPath(t *testing.T) {
	homePath, _ := homedir.Dir()
	dir, err := getConfigPath()
	require.Nil(t, err)
	require.Equal(t, dir, filepath.Join(homePath, configDirectory))
}
