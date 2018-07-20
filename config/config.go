package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const envPrefix = "MESG"
const configFileName = "config"

func initViperEnv() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
}
