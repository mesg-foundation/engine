package config

import (
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stvp/assert"
)

func TestCreateConfigDirectory(t *testing.T) {
	err := createConfigDirectory()
	assert.Nil(t, err)
}

func TestConfigDirectory(t *testing.T) {
	homeDirectory, _ := homedir.Dir()
	dir, err := getConfigDirectory()
	assert.Nil(t, err)
	assert.Equal(t, dir, filepath.Join(homeDirectory, configDirectory))
}
