package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/docker"
	"github.com/spf13/viper"
	"github.com/stvp/assert"
)

func TestStart(t *testing.T) {
	<-testForceAndWaitForFullStop()
	service, err := Start()
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestStartNetwork(t *testing.T) {
	Start()
	network, err := docker.FindNetwork(namespaceNetwork())
	assert.Nil(t, err)
	assert.NotNil(t, network)
}

func TestStartConfig(t *testing.T) {
	service, err := Start()
	assert.Nil(t, err)
	// Make sure that the config directory is passed in parameter to write on the same folder
	assert.Equal(t, service.Spec.TaskTemplate.ContainerSpec.Env[0], "MESG.PATH=/mesg")
	// Ensure that the port is shared
	assert.Equal(t, service.Spec.EndpointSpec.Ports[0].PublishedPort, uint32(50052))
	assert.Equal(t, service.Spec.EndpointSpec.Ports[0].TargetPort, uint32(50052))
	// Ensure that the docker socket is shared in the daemon
	assert.Equal(t, service.Spec.TaskTemplate.ContainerSpec.Mounts[0].Source, "/var/run/docker.sock")
	assert.Equal(t, service.Spec.TaskTemplate.ContainerSpec.Mounts[0].Target, "/var/run/docker.sock")
	// Ensure that the host users folder is sync with the daemon
	assert.Equal(t, service.Spec.TaskTemplate.ContainerSpec.Mounts[1].Source, viper.GetString(config.MESGPath))
	assert.Equal(t, service.Spec.TaskTemplate.ContainerSpec.Mounts[1].Target, "/mesg")
}
