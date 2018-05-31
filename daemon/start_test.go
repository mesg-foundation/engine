package daemon

import (
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
	"github.com/stvp/assert"
)

// startForTest starts a dummy daemon service
func startForTest() {
	running, err := IsRunning()
	if err != nil {
		panic(err)
	}
	if running == true {
		return
	}
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		panic(err)
	}
	_, err = container.StartService(container.ServiceOptions{
		Namespace:  []string{name},
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
	spec := serviceSpec("networkID")
	assert.NotNil(t, spec)
	// Make sure that the config directory is passed in parameter to write on the same folder
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Env[0], "MESG.PATH=/mesg")
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Env[1], "API.SERVICE.SOCKETPATH="+filepath.Join(viper.GetString(config.MESGPath), "server.sock"))
	// Ensure that the port is shared
	assert.Equal(t, spec.EndpointSpec.Ports[0].PublishedPort, uint32(50052))
	assert.Equal(t, spec.EndpointSpec.Ports[0].TargetPort, uint32(50052))
	// Ensure that the docker socket is shared in the daemon
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Mounts[0].Source, "/var/run/docker.sock")
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Mounts[0].Target, "/var/run/docker.sock")
	// Ensure that the host users folder is sync with the daemon
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Mounts[1].Source, viper.GetString(config.MESGPath))
	assert.Equal(t, spec.TaskTemplate.ContainerSpec.Mounts[1].Target, "/mesg")
}
