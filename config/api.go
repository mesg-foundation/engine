package config

import (
	"os"
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

// All the configuration keys.
const (
	APIServerAddress = "Api.Server.Address"
	APIClientTarget  = "Api.Client.Target"
	LogFormat        = "Log.Format"
	LogLevel         = "Log.Level"
	MESGPath         = "MESG.Path"
	CoreImage        = "Core.Image"
)

func setAPIDefault() {
	configPath, _ := getConfigPath()

	viper.SetDefault(MESGPath, configPath)

	viper.SetDefault(APIServerAddress, ":50052")
	os.MkdirAll("/mesg", os.ModePerm)

	viper.SetDefault(APIClientTarget, viper.GetString(APIServerAddress))

	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
