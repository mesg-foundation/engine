package config

import (
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/sirupsen/logrus"
)

// Path to a dedicated directory for Core
const Path = "/mesg"

// APIPort is the port of the GRPC API
func APIPort() *Config {
	return new(newViperSetting("API.Port", "50052"))
}

// APIAddress is the ip address of the GRPC API
func APIAddress() *Config {
	return new(newViperSetting("API.Address", ""))
}

// LogFormat is the log's format. Can be text or JSON.
func LogFormat() *Config {
	return new(newViperSetting("Log.Format", "text"), withAllowedValues("text", "json"))
}

// LogLevel is the minimum log's level to output.
func LogLevel() *Config {
	validation := func(value string) error {
		if _, err := logrus.ParseLevel(value); err != nil {
			return fmt.Errorf("config: %s is not valid log level", value)
		}
		return nil
	}
	return new(newViperSetting("Log.Level", "info"), withValidation(validation))
}

// CoreImage is the port of the GRPC API
func CoreImage() *Config {
	coreTag := strings.Split(version.Version, " ")
	return new(newViperSetting("Core.Image", "mesg/core:"+coreTag[0]))
}
