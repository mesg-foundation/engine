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

func getEnv() (map[string]string, error) {
	logFormat := config.LogFormat()
	logFormatValue, err := logFormat.GetValue()
	if err != nil {
		return map[string]string{}, err
	}
	logLevel := config.LogLevel()
	logLevelValue, err := logLevel.GetValue()
	if err != nil {
		return map[string]string{}, err
	}
	return map[string]string{
		logFormat.GetEnvKey(): logFormatValue,
		logLevel.GetEnvKey():  logLevelValue,
	}, nil
}

func getPorts() ([]container.Port, error) {
	portValue, err := config.APIPort().GetValue()
	if err != nil {
		return []container.Port{}, err
	}
	port, err := strconv.Atoi(portValue)
	if err != nil {
		return []container.Port{}, err
	}
	return []container.Port{
		{
			Target:    uint32(port),
			Published: uint32(port),
		},
	}, nil
}

func serviceSpec() (spec container.ServiceOptions, err error) {
	sharedNetworkID, err := defaultContainer.SharedNetworkID()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	coreImage, err := config.CoreImage().GetValue()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	env, err := getEnv()
	if err != nil {
		return container.ServiceOptions{}, err
	}
	ports, err := getPorts()
	if err != nil {
		return container.ServiceOptions{}, err
	}

	return container.ServiceOptions{
		Namespace: Namespace(),
		Image:     coreImage,
		Env:       container.MapToEnv(env),
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
		Ports:      ports,
		NetworksID: []string{sharedNetworkID},
	}, nil
}
