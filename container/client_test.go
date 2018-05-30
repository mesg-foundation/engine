package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestClient(t *testing.T) {
	client, err := Client()
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestClientIsTheSame(t *testing.T) {
	client, _ := Client()
	client2, err := Client()
	assert.Nil(t, err)
	assert.NotNil(t, client2)
	assert.Equal(t, client, client2)
}

func TestClientNotIsTheSame(t *testing.T) {
	client, _ := Client()
	resetClient()
	client2, err := Client()
	assert.Nil(t, err)
	assert.NotNil(t, client2)
	assert.NotEqual(t, client, client2)
}

// TODO: this tests break other tests on my machine
// func TestCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	ID, err := createSwarm(dockerClient)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, ID, "")
// }

// func TestClientWithCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	client, err := Client()
// 	assert.Nil(t, err)
// 	assert.NotNil(t, client)
// }

// func leaveSwarm() {
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	dockerClient.LeaveSwarm(godocker.LeaveSwarmOptions{
// 		Context: context.Background(),
// 		Force:   true,
// 	})
// }
