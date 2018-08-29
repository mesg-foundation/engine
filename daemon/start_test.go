package daemon

import (
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
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

// func TestStart(t *testing.T) {
// 	<-testForceAndWaitForFullStop()
// 	service, err := Start()
// 	require.Nil(t, err)
// 	require.NotNil(t, service)
// }

func contains(list []string, item string) bool {
	for _, itemInList := range list {
		if itemInList == item {
			return true
		}
	}
	return false
}

func TestStartConfig(t *testing.T) {
	spec, err := serviceSpec()
	require.Nil(t, err)
	// Make sure that the config directory is passed in parameter to write on the same folder
	require.True(t, contains(spec.Env, "MESG_MESG_PATH=/mesg"))
	require.True(t, contains(spec.Env, "MESG_API_SERVICE_SOCKETPATH="+filepath.Join(viper.GetString(config.MESGPath), "server.sock")))
	require.True(t, contains(spec.Env, "MESG_SERVICE_PATH_HOST="+filepath.Join(viper.GetString(config.MESGPath), "services")))
	// Ensure that the port is shared
	require.Equal(t, spec.Ports[0].Published, uint32(50052))
	require.Equal(t, spec.Ports[0].Target, uint32(50052))
	// Ensure that the docker socket is shared in the core
	require.Equal(t, spec.Mounts[0].Source, dockerSocket)
	require.Equal(t, spec.Mounts[0].Target, dockerSocket)
	// Ensure that the host users folder is sync with the core
	require.Equal(t, spec.Mounts[1].Source, viper.GetString(config.MESGPath))
	require.Equal(t, spec.Mounts[1].Target, "/mesg")
}
