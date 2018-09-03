package config

import (
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

// All the configuration keys.
const (
	APIPort    = "API.Port"    // The port of the GRPC API
	APIAddress = "API.Address" // The ip address of the GRPC API
	LogFormat  = "Log.Format"  // The log's format. Can be text or JSON
	LogLevel   = "Log.Level"   // The minimum log's level to output
	CoreImage  = "Core.Image"  // The Core's image to use
)

func setAPIDefault() {
	viper.SetDefault(Path, "/mesg")

	viper.SetDefault(APIPort, "50052")
	viper.SetDefault(APIAddress, "")

	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])
}
