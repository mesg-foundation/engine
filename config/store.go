package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mesg-foundation/core/version"
	"github.com/sirupsen/logrus"
)

// Path to a dedicated directory for Core
// TODO: Path should be reverted to a const when the package database is renovated
var Path = "/mesg"

var (
	storeInstance *store
	storeOnce     sync.Once
)

// store contains all necessary configs
type store struct {
	apiPort    *Config
	apiAddress *Config
	logFormat  *Config
	logLevel   *Config
	coreImage  *Config
}

// getStore return a singleton of store
func getStore() *store {
	storeOnce.Do(func() {
		viperEngine := getViperEngine()
		storeInstance = &store{
			apiPort:    new("API.Port", "50052", viperEngine),
			apiAddress: new("API.Address", "", viperEngine),
			logFormat:  new("Log.Format", "text", viperEngine, withAllowedValues("text", "json")),
			logLevel: new("Log.Level", "info", getViperEngine(), withValidation(func(value string) error {
				if _, err := logrus.ParseLevel(value); err != nil {
					return fmt.Errorf("Value %q is not a valid log level", value)
				}
				return nil
			})),
			coreImage: new("Core.Image", "mesg/core:"+strings.Split(version.Version, " ")[0], viperEngine),
		}
	})
	return storeInstance
}

// APIPort is the port of the GRPC API
func APIPort() *Config {
	return getStore().apiPort
}

// APIAddress is the ip address of the GRPC API
func APIAddress() *Config {
	return getStore().apiAddress
}

// LogFormat is the log's format. Can be text or JSON.
func LogFormat() *Config {
	return getStore().logFormat
}

// LogLevel is the minimum log's level to output.
func LogLevel() *Config {
	return getStore().logLevel
}

// CoreImage is the port of the GRPC API
func CoreImage() *Config {
	return getStore().coreImage
}
