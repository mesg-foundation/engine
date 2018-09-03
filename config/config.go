package config

import (
	"strings"

	"github.com/mesg-foundation/core/version"
)

// Path to a dedicated directory for Core
const Path = "/mesg"

// APIPort is the port of the GRPC API
func APIPort() Setting {
	return newViperSetting("API.Port", "50052")
}

// APIAddress is the ip address of the GRPC API
func APIAddress() Setting {
	return newViperSetting("API.Address", "")
}

// LogFormat is the log's format. Can be text or JSON.
func LogFormat() Setting {
	return newViperSetting("Log.Format", "text")
}

// LogLevel is the minimum log's level to output.
func LogLevel() Setting {
	return newViperSetting("Log.Level", "info")
}

// CoreImage is the port of the GRPC API
func CoreImage() Setting {
	coreTag := strings.Split(version.Version, " ")
	return newViperSetting("Core.Image", "mesg/core:"+coreTag[0])
}

// func validateConfig() {
// 	format := viper.GetString(LogFormat)
// 	if format != "text" && format != "json" {
// 		fmt.Fprintf(os.Stderr, "config: %s is not valid log format", format)
// 		os.Exit(1)
// 	}

// 	level := viper.GetString(LogLevel)
// 	if _, err := logrus.ParseLevel(level); err != nil {
// 		fmt.Fprintf(os.Stderr, "config: %s is not valid log level", level)
// 		os.Exit(1)
// 	}
// }
