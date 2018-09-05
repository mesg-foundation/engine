package config

import (
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/version"
	"github.com/sirupsen/logrus"
)

// Path to a dedicated directory for Core
// TODO: Path should be reverted to a const when the package database is renovated
var Path = "/mesg"

// APIPort is the port of the GRPC API
func APIPort() *Config {
	return new("API.Port", "50052", getViperEngine())
}

// APIAddress is the ip address of the GRPC API
func APIAddress() *Config {
	return new("API.Address", "", getViperEngine())
}

// LogFormat is the log's format. Can be text or JSON.
func LogFormat() *Config {
	return new("Log.Format", "text", getViperEngine(), withAllowedValues("text", "json"))
}

// LogLevel is the minimum log's level to output.
func LogLevel() *Config {
	validation := func(value string) error {
		if _, err := logrus.ParseLevel(value); err != nil {
			return fmt.Errorf("Value %q is not a valid log level", value)
		}
		return nil
	}
	return new("Log.Level", "info", getViperEngine(), withValidation(validation))
}

// CoreImage is the port of the GRPC API
func CoreImage() *Config {
	coreTag := strings.Split(version.Version, " ")
	return new("Core.Image", "mesg/core:"+coreTag[0], getViperEngine())
}
