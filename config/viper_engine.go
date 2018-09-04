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
	viperEngineInstance *viperEngine
	viperEngineOnce     sync.Once
)

type viperEngine struct {
	viper *viper.Viper
}

// getViperEngine returns the viperEngine and init it if needed
func getViperEngine() *viperEngine {
	viperEngineOnce.Do(func() {
		viperEngineInstance = &viperEngine{
			viper: viper.New(),
		}
		viperEngineInstance.readEnv()
		viperEngineInstance.readConfigFile()
	})
	return viperEngineInstance
}

// readEnv populates viper from the env variable
func (v *viperEngine) readEnv() {
	v.viper.SetEnvPrefix(envPrefix)
	v.viper.AutomaticEnv()
	v.viper.SetEnvKeyReplacer(strings.NewReplacer(defaultSeparator, envSeparator))
}

// readConfigFile populates viper from the config file
func (v *viperEngine) readConfigFile() {
	v.viper.SetConfigName(configFileName)
	v.viper.AddConfigPath("$HOME") // for user home path
	v.viper.AddConfigPath(".")     // for current path
	if v.viper.ReadInConfig() == nil {
		fmt.Println("Using config file:", v.viper.ConfigFileUsed())
	}
}

func (v *viperEngine) setDefaultValue(key string, defaultValue string) error {
	v.viper.SetDefault(key, defaultValue)
	return nil
}

func (v *viperEngine) setValue(key string, value string) error {
	v.viper.Set(key, value)
	return nil
}

func (v *viperEngine) getValue(key string) (string, error) {
	value := v.viper.GetString(key)
	return value, nil
}

func (v *viperEngine) getEnvKey(key string) string {
	replacer := strings.NewReplacer(defaultSeparator, envSeparator)
	return envPrefix + envSeparator + replacer.Replace(strings.ToUpper(key))
}
