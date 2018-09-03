package config

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func resetViperInstance() {
	viperInstance = nil
	viperOnce = sync.Once{}
}

func TestGetViper(t *testing.T) {
	resetViperInstance()
	viper := getViper()
	require.NotNil(t, viper)
	viper2 := getViper()
	require.Equal(t, viper, viper2)
}

func TestViperReadEnv(t *testing.T) {
	resetViperInstance()
	os.Setenv(envPrefix+envSeparator+"TESTVIPERREADENV", "envValue")
	viper := getViper()
	require.Equal(t, "envValue", viper.GetString("TestViperReadEnv"))
}

func TestViperReadConfigFile(t *testing.T) {
	resetViperInstance()
	path := "./" + configFileName + ".yml"
	content := []byte("test: hello\n")
	ioutil.WriteFile(path, content, 0644)
	defer os.Remove(path)
	viper := getViper()
	require.Equal(t, "hello", viper.GetString("test"))
}
