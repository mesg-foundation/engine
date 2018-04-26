package service

import (
	"context"
	"os"
	"testing"

	"github.com/fsouza/go-dockerclient"

	"github.com/stvp/assert"
)

func TestDockerCliSingleton(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	cli, err := dockerCli()
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestCreateDockerCli(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	_, err := createDockerCli()
	assert.Nil(t, err)
}

func TestCreateDockerCliWithSwarm(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	cli, _ := dockerCli()
	cli.LeaveSwarm(docker.LeaveSwarmOptions{
		Context: context.Background(),
		Force:   true,
	})
	dockerCliInstance = nil
	_, err := createDockerCli()
	assert.Nil(t, err)
}

func TestCreateSwarm(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	cli, _ := dockerCli()
	cli.LeaveSwarm(docker.LeaveSwarmOptions{
		Context: context.Background(),
		Force:   true,
	})
	ID, err := createSwarm(cli)
	assert.Nil(t, err)
	assert.NotEqual(t, ID, "")
}
