package container

import (
	"context"
	"testing"

	docker "github.com/docker/docker/client"
	"github.com/stvp/assert"
)

func removeSharedNetworkIfExist(client *docker.Client) (err error) {
	_, err = sharedNetwork(client)
	if docker.IsErrNotFound(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}
	err = client.NetworkRemove(context.Background(), Namespace(sharedNetworkNamespace))
	return
}

func TestCreateSharedNetworkIfNeeded(t *testing.T) {
	client, _ := createClient()
	err := removeSharedNetworkIfExist(client)
	assert.Nil(t, err)
	err = createSharedNetworkIfNeeded(client)
	assert.Nil(t, err)
}

func TestSharedNetwork(t *testing.T) {
	client, _ := Client()
	network, err := sharedNetwork(client)
	assert.Nil(t, err)
	assert.NotEqual(t, "", network.ID)
}

func TestSharedNetworkID(t *testing.T) {
	networkID, err := SharedNetworkID()
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
}
