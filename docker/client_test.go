package docker

import (
	"context"
	"testing"

	godocker "github.com/fsouza/go-dockerclient"
	"github.com/stvp/assert"
)

func TestDockerCliSingleton(t *testing.T) {
	cli, err := Client()
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestCreateDockerCli(t *testing.T) {
	_, err := createClient()
	assert.Nil(t, err)
}

func TestCreateDockerCliWithSwarm(t *testing.T) {
	cli, _ := Client()
	cli.LeaveSwarm(godocker.LeaveSwarmOptions{
		Context: context.Background(),
		Force:   true,
	})
	dockerCliInstance = nil
	_, err := createClient()
	assert.Nil(t, err)
}

func TestCreateSwarm(t *testing.T) {
	cli, _ := Client()
	cli.LeaveSwarm(godocker.LeaveSwarmOptions{
		Context: context.Background(),
		Force:   true,
	})
	ID, err := createSwarm(cli)
	assert.Nil(t, err)
	assert.NotEqual(t, ID, "")
}
