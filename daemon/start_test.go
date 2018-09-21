package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/stretchr/testify/require"
)

// startForTest starts a dummy MESG Core service
func startForTest() {
	status, err := Status()
	if err != nil {
		panic(err)
	}
	if status == container.RUNNING {
		return
	}
	sharedNetworkID, err := defaultContainer.SharedNetworkID()
	if err != nil {
		panic(err)
	}
	_, err = defaultContainer.StartService(container.ServiceOptions{
		Namespace:  Namespace(),
		Image:      "http-server",
		NetworksID: []string{sharedNetworkID},
	})
	if err != nil {
		panic(err)
	}
	return
}

func TestStartConfig(t *testing.T) {
	c, _ := config.Global()
	spec, err := serviceSpec()
	require.NoError(t, err)
	// Make sure that the config directory is passed in parameter to write on the same folder
	require.Contains(t, spec.Env, "MESG_LOG_LEVEL=info")
	require.Contains(t, spec.Env, "MESG_LOG_FORMAT=text")
	require.Contains(t, spec.Env, "MESG_CORE_PATH="+c.Docker.Core.Path)
	// Ensure that the port is shared
	_, port, _ := xnet.SplitHostPort(c.Server.Address)
	require.Equal(t, spec.Ports[0].Published, uint32(port))
	require.Equal(t, spec.Ports[0].Target, uint32(port))
	// Ensure that the docker socket is shared in the core
	require.Equal(t, spec.Mounts[0].Source, c.Docker.Socket)
	require.Equal(t, spec.Mounts[0].Target, c.Docker.Socket)
	require.True(t, spec.Mounts[0].Bind)
	// Ensure that the host users folder is sync with the core
	require.Equal(t, spec.Mounts[1].Source, c.Core.Path)
	require.Equal(t, spec.Mounts[1].Target, c.Docker.Core.Path)
	require.True(t, spec.Mounts[1].Bind)
}
