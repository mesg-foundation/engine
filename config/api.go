package config

import (
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

// All the configuration keys.
const (
	LogFormat        = "Log.Format"
	LogLevel         = "Log.Level"
	CoreImage        = "Core.Image"
	APIGRPCPort     = "API.GRPC.Port"
	APIGRPCAddresss = "API.GRPC.Address"
	Path            = "Path" // The path to a dedicated directory for Core
)

func setAPIDefault() {
	viper.SetDefault(Path, "/mesg")

	viper.SetDefault(APIGRPCPort, 50052)
	viper.SetDefault(APIGRPCAddresss, ":50052")

	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
