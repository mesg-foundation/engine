package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	envPrefix        = "MESG"
	envSeparator     = "_"
	defaultSeparator = "."
	configFileName   = ".mesg"
)

var (
	viperInstance *viper.Viper
	viperOnce     sync.Once
)

// getViper returns the viperInstance and init it if needed
func getViper() *viper.Viper {
	viperOnce.Do(func() {
		viperInstance = viper.New()
		readEnv(viperInstance)
		readConfigFile(viperInstance)
	})
	return viperInstance
}

// readEnv populates viper from the env variable
func readEnv(viper *viper.Viper) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(defaultSeparator, envSeparator))
}

// readConfigFile populates viper from the config file
func readConfigFile(viper *viper.Viper) {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("$HOME") // for user home path
	viper.AddConfigPath(".")     // for current path
	if viper.ReadInConfig() == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
