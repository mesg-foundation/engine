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
	return c.client.NetworkRemove(context.Background(), Namespace(sharedNetworkNamespace))
}

func TestIntegrationCreateSharedNetworkIfNeeded(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	err = removeSharedNetworkIfExist(c)
	require.NoError(t, err)
	err = c.createSharedNetworkIfNeeded()
	require.NoError(t, err)
}

func TestIntegrationSharedNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	network, err := c.sharedNetwork()
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationSharedNetworkID(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	networkID, err := c.SharedNetworkID()
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}
