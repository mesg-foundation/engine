package config

import (
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stvp/assert"
)

func TestCreateConfigPath(t *testing.T) {
	err := createConfigPath()
	assert.Nil(t, err)
}

func TestConfigPath(t *testing.T) {
	homePath, _ := homedir.Dir()
	dir, err := getConfigPath()
	assert.Nil(t, err)
	assert.Equal(t, dir, filepath.Join(homePath, configDirectory))
}
