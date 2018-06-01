package daemon

import (
	"path/filepath"
	"time"

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
	serviceID, err = container.StartService(serviceSpec(sharedNetworkID))
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(Namespace(), container.RUNNING, time.Minute)
	return
}

func serviceSpec(networkID string) container.ServiceOptions {
	return container.ServiceOptions{
		Namespace: Namespace(),
		Image:     viper.GetString(config.DaemonImage),
		Env: []string{
			"MESG.PATH=/mesg",
			"API.SERVICE.SOCKETPATH=" + filepath.Join(viper.GetString(config.MESGPath), "server.sock"),
			"SERVICE.PATH.HOST=" + filepath.Join(viper.GetString(config.MESGPath), "services"),
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
