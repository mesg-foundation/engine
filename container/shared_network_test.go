package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
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
	require.Nil(t, err)
	require.Equal(t, id, network.ID)

	li := <-dt.LastNetworkInspect()
	require.Equal(t, Namespace([]string{}), li.Network)
	require.Equal(t, types.NetworkInspectOptions{}, li.Options)
}

func TestCreateSharedNetworkIfNeeded(t *testing.T) {
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{}, nil)

	require.Nil(t, c.createSharedNetworkIfNeeded())

	lc := <-dt.LastNetworkCreate()
	require.Equal(t, Namespace([]string{}), lc.Name)
	require.Equal(t, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": Namespace([]string{}),
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

	require.Nil(t, c.createSharedNetworkIfNeeded())

	select {
	case <-dt.LastNetworkCreate():
		t.Error("should not create network")
	default:
	}
}

func TestSharedNetworkID(t *testing.T) {
	id := "1"
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	network, err := c.SharedNetworkID()
	require.Nil(t, err)
	require.Equal(t, network, id)

	li := <-dt.LastNetworkInspect()
	require.Equal(t, Namespace([]string{}), li.Network)
	require.Equal(t, types.NetworkInspectOptions{}, li.Options)
}
