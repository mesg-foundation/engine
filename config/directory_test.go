package config

import (
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultPath(t *testing.T) {
	homePath, _ := homedir.Dir()
	dir := getDefaultPath()
	require.Equal(t, dir, filepath.Join(homePath, defaultDirectory))
}

func TestCreatePath(t *testing.T) {
	err := createPath()
	require.Nil(t, err)
}

func TestCreateServicesPath(t *testing.T) {
	err := createServicesPath()
	require.Nil(t, err)
}
