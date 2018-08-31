package config

import (
	"path/filepath"
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

// All the API configuration keys.
const (
	APIServerAddress  = "Api.Server.Address"
	APIClientTarget   = "Api.Client.Target"
	ServicePathHost   = "Service.Path.Host"
	ServicePathDocker = "Service.Path.Docker"
	CoreImage         = "Core.Image"
)

func setAPIDefault() {
	viper.SetDefault(APIServerAddress, ":50052")

	viper.SetDefault(APIClientTarget, viper.GetString(APIServerAddress))

	viper.SetDefault(ServicePathHost, filepath.Join(viper.GetString(MESGPath), "services"))
	viper.SetDefault(ServicePathDocker, filepath.Join("/mesg", "services"))

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
