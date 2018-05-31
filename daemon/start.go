package daemon

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start the docker daemon
func Start() (serviceID string, err error) {
	running, err := IsRunning()
	if err != nil {
		return
	}
	if running == true {
		return
	}
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		return
	}
	return container.StartService(serviceSpec(sharedNetworkID))
}

func serviceSpec(networkID string) container.ServiceOptions {
	return container.ServiceOptions{
		Namespace: []string{name},
		Image:     viper.GetString(config.DaemonImage),
		Env: []string{
			"MESG.PATH=/mesg",
			"API.SERVICE.SOCKETPATH=" + filepath.Join(viper.GetString(config.MESGPath), "server.sock"),
		},
		Mounts: []container.Mount{
			container.Mount{
				Source: dockerSocket,
				Target: dockerSocket,
			},
			container.Mount{
				Source: viper.GetString(config.MESGPath),
				Target: "/mesg",
			},
		},
		Ports: []container.Port{
			container.Port{
				Target:    50052,
				Published: 50052,
			},
		},
		NetworksID: []string{networkID},
	}
}
