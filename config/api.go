package config

import (
	"os"
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
	MESGPath          = "MESG.Path"
	CoreImage         = "Core.Image"
)

func setAPIDefault() {
	configPath, _ := getConfigPath()

	viper.SetDefault(MESGPath, configPath)

	viper.SetDefault(APIServerAddress, ":50052")
	os.MkdirAll("/mesg", os.ModePerm)

	viper.SetDefault(APIClientTarget, viper.GetString(APIServerAddress))

	viper.SetDefault(ServicePathHost, filepath.Join(viper.GetString(MESGPath), "services"))
	viper.SetDefault(ServicePathDocker, filepath.Join("/mesg", "services"))
	os.MkdirAll(viper.GetString(ServicePathDocker), os.ModePerm)

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
