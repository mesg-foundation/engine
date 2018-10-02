// +build integration

package container

import (
	"context"
	"testing"

	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

func removeSharedNetworkIfExist(c *Container) error {
	if _, err := c.sharedNetwork(); err != nil {
		if docker.IsErrNotFound(err) {
			return nil
		}
		return err
	}
	return c.client.NetworkRemove(context.Background(), c.Namespace([]string{}))
}

func TestIntegrationCreateSharedNetworkIfNeeded(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	err = removeSharedNetworkIfExist(c)
	require.Nil(t, err)
	err = c.createSharedNetworkIfNeeded()
	require.Nil(t, err)
}

func TestIntegrationSharedNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	network, err := c.sharedNetwork()
	require.Nil(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationSharedNetworkID(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	networkID, err := c.SharedNetworkID()
	require.Nil(t, err)
	require.NotEqual(t, "", networkID)
}
