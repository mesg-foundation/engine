package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateDockerCli(t *testing.T) {
	_, err := createDockerCli("unix:///var/run/docker.sock")
	assert.Nil(t, err)
}
