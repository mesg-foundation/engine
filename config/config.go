package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	envPrefix        = "MESG"
	envSeparator     = "_"
	defaultSeparator = "."
	configFileName   = "config"
)

// ToEnv transforms a config key to a env key.
func ToEnv(key string) string {
	replacer := strings.NewReplacer(defaultSeparator, envSeparator)
	return envPrefix + envSeparator + replacer.Replace(strings.ToUpper(key))
}

func readConfigFromEnv() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(defaultSeparator, envSeparator))
}

func readConfigFromFile() {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(viper.GetString(MESGPath))
	if viper.ReadInConfig() == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	// The order of the following functions is critical. Do not change it.
	// Set the default app directory
	setDirectoryDefault()

	// Read the app directory from env if set
	readConfigFromEnv()

	// Create the required directories if needed
	if err := createPath(); err != nil {
		panic(err)
	}
	if err := createServicesPath(); err != nil {
		panic(err)
	}

	// Read the config file from the app directory if exist
	readConfigFromFile()

	setCoreDefault()
	setLogDefault()

	validateLog()
}
