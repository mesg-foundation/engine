package daemon

import (
	"strconv"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
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
	portValue, err := config.APIPort().GetValue()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	port, err := strconv.Atoi(portValue)
	if err != nil {
		return container.ServiceOptions{}, err
	}
	coreImage, err := config.CoreImage().GetValue()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	logFormat := config.LogFormat()
	logFormatValue, err := logFormat.GetValue()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	logLevel := config.LogLevel()
	logLevelValue, err := logLevel.GetValue()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	return container.ServiceOptions{
		Namespace: Namespace(),
		Image:     coreImage,
		Env: container.MapToEnv(map[string]string{
			logFormat.GetEnvKey(): logFormatValue,
			logLevel.GetEnvKey():  logLevelValue,
		}),
		Mounts: []container.Mount{
			{
				Source: dockerSocket,
				Target: dockerSocket,
				Bind:   true,
			},
			{
				Source: volume,
				Target: config.Path,
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
