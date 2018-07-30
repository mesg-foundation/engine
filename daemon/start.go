package daemon

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start the docker core
func Start() (serviceID string, err error) {
	status, err := Status()
	if err != nil || status == container.RUNNING {
		return
	}
	spec, err := serviceSpec()
	if err != nil {
		return
	}
	serviceID, err = container.StartService(spec)
	return
}

func serviceSpec() (spec container.ServiceOptions, err error) {
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		return
	}
	spec = container.ServiceOptions{
		Namespace: Namespace(),
		Image:     viper.GetString(config.CoreImage),
		Env: container.MapToEnv(map[string]string{
			config.ToEnv(config.MESGPath):             "/mesg",
			config.ToEnv(config.APIServiceSocketPath): filepath.Join(viper.GetString(config.MESGPath), "server.sock"),
			config.ToEnv(config.ServicePathHost):      filepath.Join(viper.GetString(config.MESGPath), "services"),
		}),
		Mounts: []container.Mount{
			{
				Source: dockerSocket,
				Target: dockerSocket,
			},
			{
				Source: viper.GetString(config.MESGPath),
				Target: "/mesg",
			},
		},
		Ports: []container.Port{
			{
				Target:    50052,
				Published: 50052,
			},
		},
		NetworksID: []string{sharedNetworkID},
	}
	return
}
