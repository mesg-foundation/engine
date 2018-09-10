package config

import (
	"strings"
	"sync"

	"github.com/mesg-foundation/core/version"
	"github.com/sirupsen/logrus"
)

// Path to a dedicated directory for Core
// TODO: Path should be reverted to a const when the package database is renovated
var Path = "/mesg"

var (
	configInstance *config
	configOnce     sync.Once
)

// config contains all necessary configs
type config struct {
	apiPort    *Entry
	apiAddress *Entry
	logFormat  *Entry
	logLevel   *Entry
	coreImage  *Entry
}

// getConfig return a singleton of config
func getConfig() *config {
	configOnce.Do(func() {
		viperEngine := getViperEngine()
		configInstance = &config{
			apiPort:    newEntry("API.Port", "50052", viperEngine),
			apiAddress: newEntry("API.Address", "", viperEngine),
			logFormat:  newEntry("Log.Format", "text", viperEngine, withAllowedValues("text", "json")),
			logLevel: newEntry("Log.Level", "info", viperEngine, withValidation(func(value string) error {
				_, err := logrus.ParseLevel(value)
				return err
			})),
			coreImage: newEntry("Core.Image", "mesg/core:"+strings.Split(version.Version, " ")[0], viperEngine),
		}
	})
	return configInstance
}

// APIPort is the port of the GRPC API
func APIPort() *Entry {
	return getConfig().apiPort
}

// APIAddress is the ip address of the GRPC API
func APIAddress() *Entry {
	return getConfig().apiAddress
}

// LogFormat is the log's format. Can be text or JSON.
func LogFormat() *Entry {
	return getConfig().logFormat
}

// LogLevel is the minimum log's level to output.
func LogLevel() *Entry {
	return getConfig().logLevel
}

// CoreImage is the port of the GRPC API
func CoreImage() *Entry {
	return getConfig().coreImage
}
