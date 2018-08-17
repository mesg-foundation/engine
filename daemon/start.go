package daemon

import (
	"net"
	"path/filepath"
	"strconv"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start starts the docker core.
func Start() (serviceID string, err error) {
	status, err := Status()
	if err != nil || status == container.RUNNING {
		return "", err
	}
	spec, err := serviceSpec()
	if err != nil {
		return "", err
	}
	return defaultContainer.StartService(spec)
}

func serviceSpec() (spec container.ServiceOptions, err error) {
	sharedNetworkID, err := defaultContainer.SharedNetworkID()
	if err != nil {
		return container.ServiceOptions{}, err
	}

	_, portStr, err := net.SplitHostPort(viper.GetString(config.APIServerAddress))
	if err != nil {
		return container.ServiceOptions{}, err
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		return container.ServiceOptions{}, err
	}

	return container.ServiceOptions{
		Namespace: Namespace(),
		Image:     viper.GetString(config.CoreImage),
		Env: container.MapToEnv(map[string]string{
			config.ToEnv(config.MESGPath):             "/mesg",
			config.ToEnv(config.LogFormat):            viper.GetString(config.LogFormat),
			config.ToEnv(config.LogLevel):             viper.GetString(config.LogLevel),
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
				Target:    uint32(port),
				Published: uint32(port),
			},
		},
		NetworksID: []string{sharedNetworkID},
	}, nil
}
