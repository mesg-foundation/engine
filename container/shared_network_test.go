package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestSharedNetwork(t *testing.T) {
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	network, err := c.sharedNetwork()
	assert.Nil(t, err)
	assert.Equal(t, id, network.ID)

	li := <-dt.LastNetworkInspect()
	assert.Equal(t, Namespace(sharedNetworkNamespace), li.Network)
	assert.Equal(t, types.NetworkInspectOptions{}, li.Options)
}

func TestCreateSharedNetworkIfNeeded(t *testing.T) {
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{}, nil)

	assert.Nil(t, c.createSharedNetworkIfNeeded())

	lc := <-dt.LastNetworkCreate()
	assert.Equal(t, Namespace(sharedNetworkNamespace), lc.Name)
	assert.Equal(t, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": Namespace(sharedNetworkNamespace),
		},
	}, lc.Options)
}

func TestCreateSharedNetworkIfNeededExists(t *testing.T) {
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	assert.Nil(t, c.createSharedNetworkIfNeeded())

	select {
	case <-dt.LastNetworkCreate():
		t.Error("should not create network")
	default:
	}
}

func TestIntegrationSharedNetworkID(t *testing.T) {
	id := "1"
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	network, err := c.SharedNetworkID()
	assert.Nil(t, err)
	assert.Equal(t, network, id)

	li := <-dt.LastNetworkInspect()
	assert.Equal(t, Namespace(sharedNetworkNamespace), li.Network)
	assert.Equal(t, types.NetworkInspectOptions{}, li.Options)
}
