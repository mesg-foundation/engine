package daemon

import (
	"runtime"
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

	//windows hack - TODO: Put into utility package
	mesgHomePath := viper.GetString(config.MESGPath)
	if runtime.GOOS == "windows" {
		mesgHomePath = "/c" + filepath.ToSlash(mesgHomePath[2:])
	}

	spec = container.ServiceOptions{
		Namespace: Namespace(),
		Image:     viper.GetString(config.CoreImage),
		Env: []string{
			"MESG.PATH=/mesg",
			//"API.SERVICE.SOCKETPATH=/c/Users/ahonegger/.mesg/server.sock",
			//"SERVICE.PATH.HOST=/c/Users/ahonegger/.mesg/services",
		},
		Mounts: []container.Mount{
			{
				Source: getDockerSocket(),
				Target: dockerSocket,
			},
			{
				Source: mesgHomePath,
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

func getDockerSocket() (dc string) {
	dc = dockerSocket
	if runtime.GOOS == "windows" {
		dc = "/" + dockerSocket
	}
	return
}