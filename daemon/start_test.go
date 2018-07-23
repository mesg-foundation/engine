package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
	"github.com/stvp/assert"
	"runtime"
	"path/filepath"
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

func TestStartConfig(t *testing.T) {
	spec, err := serviceSpec()
	assert.Nil(t, err)
	// Make sure that the config directory is passed in parameter to write on the same folder
	assert.Equal(t, spec.Env[0], "MESG.PATH=/mesg")
	//assert.Equal(t, spec.Env[1], "API.SERVICE.SOCKETPATH="+filepath.Join(viper.GetString(config.MESGPath), "server.sock"))
	//assert.Equal(t, spec.Env[2], "SERVICE.PATH.HOST="+filepath.Join(viper.GetString(config.MESGPath), "services"))
	// Ensure that the port is shared
	assert.Equal(t, spec.Ports[0].Published, uint32(50052))
	assert.Equal(t, spec.Ports[0].Target, uint32(50052))

	// Ensure that the docker socket is shared in the core
	//windows hack - TODO: Put into a utility package to avoid code duplication
	dc := dockerSocket
	if runtime.GOOS == "windows" {
		dc = "/" + dockerSocket
	}
	assert.Equal(t, spec.Mounts[0].Source, dc)
	assert.Equal(t, spec.Mounts[0].Target, "/var/run/docker.sock")

	// Ensure that the host users folder is sync with the core
	//windows hack - TODO: Put into a utility package to avoid code duplication
	mesgHomePath := viper.GetString(config.MESGPath)
	if runtime.GOOS == "windows" {
		mesgHomePath = "/c" + filepath.ToSlash(mesgHomePath[2:])
	}

	assert.Equal(t, spec.Mounts[1].Source, mesgHomePath)
	assert.Equal(t, spec.Mounts[1].Target, "/mesg")
}
