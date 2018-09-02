package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	envPrefix        = "MESG"
	envSeparator     = "_"
	defaultSeparator = "."
	configFileName   = ".mesg"
)

// ToEnv transforms a config key to a env key.
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
	viper.AddConfigPath("$HOME") // for user home path
	viper.AddConfigPath(".") // for current path
	if viper.ReadInConfig() == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func validateConfig() {
	format := viper.GetString(LogFormat)
	if format != "text" && format != "json" {
		fmt.Fprintf(os.Stderr, "config: %s is not valid log format", format)
		os.Exit(1)
	}

	level := viper.GetString(LogLevel)
	if _, err := logrus.ParseLevel(level); err != nil {
		fmt.Fprintf(os.Stderr, "config: %s is not valid log level", level)
		os.Exit(1)
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
	validateConfig()
}
