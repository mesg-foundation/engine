package daemon

import (
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
	"github.com/stvp/assert"
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
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		panic(err)
	}
	_, err = container.StartService(container.ServiceOptions{
		Namespace:  Namespace(),
		Image:      "nginx",
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
// 	assert.Nil(t, err)
// 	assert.NotNil(t, service)
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
	assert.Nil(t, err)
	// Make sure that the config directory is passed in parameter to write on the same folder
	assert.True(t, contains(spec.Env, "MESG_MESG_PATH=/mesg"))
	assert.True(t, contains(spec.Env, "MESG_API_SERVICE_SOCKETPATH="+filepath.Join(viper.GetString(config.MESGPath), "server.sock")))
	assert.True(t, contains(spec.Env, "MESG_SERVICE_PATH_HOST="+filepath.Join(viper.GetString(config.MESGPath), "services")))
	// Ensure that the port is shared
	assert.Equal(t, spec.Ports[0].Published, uint32(50052))
	assert.Equal(t, spec.Ports[0].Target, uint32(50052))
	// Ensure that the docker socket is shared in the core
	assert.Equal(t, spec.Mounts[0].Source, dockerSocket)
	assert.Equal(t, spec.Mounts[0].Target, dockerSocket)
	// Ensure that the host users folder is sync with the core
	assert.Equal(t, spec.Mounts[1].Source, viper.GetString(config.MESGPath))
	assert.Equal(t, spec.Mounts[1].Target, "/mesg")
}
