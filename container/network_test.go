package container

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	namespace := []string{"namespace"}
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkCreate(types.NetworkCreateResponse{ID: id}, nil)
	dt.ProvideNetworkInspect(types.NetworkResource{}, nil)

	networkID, err := c.CreateNetwork(namespace)
	assert.Nil(t, err)
	assert.Equal(t, id, networkID)

	li := <-dt.LastNetworkCreate()
	assert.Equal(t, Namespace(namespace), li.Name)
	assert.Equal(t, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": Namespace(namespace),
		},
	}, li.Options)
}

func TestCreateAlreadyExistingNetwork(t *testing.T) {
	namespace := []string{"namespace"}
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	networkID, err := c.CreateNetwork(namespace)
	assert.Nil(t, err)
	assert.Equal(t, id, networkID)

	li := <-dt.LastNetworkInspect()
	assert.Equal(t, Namespace(namespace), li.Network)
	assert.Equal(t, types.NetworkInspectOptions{}, li.Options)
}

func TestDeleteNetwork(t *testing.T) {
	namespace := []string{"namespace"}
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	assert.Nil(t, c.DeleteNetwork(namespace))

	li := <-dt.LastNetworkInspect()
	assert.Equal(t, Namespace(namespace), li.Network)
	assert.Equal(t, types.NetworkInspectOptions{}, li.Options)

	assert.Equal(t, id, (<-dt.LastNetworkRemove()).Network)
}

func TestDeleteNotExistingNetwork(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{}, dockertest.NotFoundErr{})

	assert.Nil(t, c.DeleteNetwork(namespace))

	select {
	case <-dt.LastNetworkRemove():
		t.Error("should not remove non existent network")
	default:
	}
}

var errNetworkDelete = errors.New("network delete")

func TestDeleteNetworkError(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{}, errNetworkDelete)

	assert.NotNil(t, c.DeleteNetwork(namespace))
}

func TestFindNetwork(t *testing.T) {
	namespace := []string{"namespace"}
	id := "id"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{ID: id}, nil)

	network, err := c.FindNetwork(namespace)
	assert.Nil(t, err)
	assert.Equal(t, id, network.ID)

	li := <-dt.LastNetworkInspect()
	assert.Equal(t, Namespace(namespace), li.Network)
	assert.Equal(t, types.NetworkInspectOptions{}, li.Options)
}

func TestFindNotExistingNetwork(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	// discard network requests made from New.
	<-dt.LastNetworkInspect()
	<-dt.LastNetworkCreate()

	dt.ProvideNetworkInspect(types.NetworkResource{}, dockertest.NotFoundErr{})

	_, err := c.FindNetwork(namespace)
	assert.Equal(t, dockertest.NotFoundErr{}, err)
}
