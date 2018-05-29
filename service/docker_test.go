package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestDockerCli(t *testing.T) {
	cli, err := dockerCli()
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestDockerCliSingleton(t *testing.T) {
	dockerCli()
	cli, err := dockerCli()
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestDockerClient(t *testing.T) {
	cli, err := DockerClient()
	assert.Nil(t, err)
	assert.NotNil(t, cli)
}

func TestCreateDockerCli(t *testing.T) {
	_, err := createDockerCli()
	assert.Nil(t, err)
}

func TestCreateNetworkIfNeeded(t *testing.T) {
	cli, _ := dockerCli()
	network, err := SharedNetwork(cli)
	if err == nil {
		cli.RemoveNetwork(network.ID)
	}
	err = createNetworkIfNeeded(cli)
	assert.Nil(t, err)
	network, err = SharedNetwork(cli)
	assert.Nil(t, err)
	assert.NotEqual(t, network.ID, "")
}

func TestResetCliInstance(t *testing.T) {
	dockerCli()
	resetCliInstance()
	assert.Nil(t, dockerCliInstance)
}

// TODO: Reactivate these tests but because swarm is destroyed
// all networks are deleted and so all other tests a failing because
// the common network is deleted too

// func TestCreateDockerCliWithSwarm(t *testing.T) {
// 	cli, _ := dockerCli()
// 	cli.LeaveSwarm(docker.LeaveSwarmOptions{
// 		Context: context.Background(),
// 		Force:   true,
// 	})
// 	resetCliInstance()
// 	_, err := createDockerCli()
// 	assert.Nil(t, err)
// }

// func TestCreateSwarm(t *testing.T) {
// 	cli, _ := dockerCli()
// 	cli.LeaveSwarm(docker.LeaveSwarmOptions{
// 		Context: context.Background(),
// 		Force:   true,
// 	})
// 	ID, err := createSwarm(cli)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, ID, "")
// }
