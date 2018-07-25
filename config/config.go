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

// ToEnv transform a config key to a env key
func ToEnv(key string) string {
	replacer := strings.NewReplacer(defaultSeparator, envSeparator)
	return envPrefix + envSeparator + replacer.Replace(strings.ToUpper(key))
}

func initViperEnv() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(defaultSeparator, envSeparator))
}

func initConfigFile() {
	viper.SetConfigName(configFileName)
	path, _ := getConfigPath()
	viper.AddConfigPath(path)
	if viper.ReadInConfig() == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	initConfigFile()
	initViperEnv()

	err := createConfigPath()
	if err != nil {
		panic(err)
	}

	setAPIDefault()
}
