package config

import (
	"path/filepath"
	"testing"

	"github.com/stvp/assert"
)

func TestHomeDirectory(t *testing.T) {
	dir, err := getHomeDirectory()
	assert.NotEqual(t, dir, "")
	assert.Nil(t, err)
}

func TestDetectHomeDirectory(t *testing.T) {
	dir, err := detectHomePath()
	assert.NotEqual(t, dir, "")
	assert.Nil(t, err)
}

func TestConfigDirectory(t *testing.T) {
	home, _ := getHomeDirectory()
	dir, err := getConfigDirectory()
	assert.Nil(t, err)
	assert.Equal(t, dir, filepath.Join(home, ".mesg"))
}

func TestAccountDirectory(t *testing.T) {
	config, _ := getConfigDirectory()
	dir, err := getAccountDirectory()
	assert.Nil(t, err)
	assert.Equal(t, dir, filepath.Join(config, "accounts"))
}
