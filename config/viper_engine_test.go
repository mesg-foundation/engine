package config

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func resetViperEngineInstance() {
	viperEngineInstance = nil
	viperEngineOnce = sync.Once{}
}

func TestGetViperEngine(t *testing.T) {
	resetViperEngineInstance()
	viperEngine := getViperEngine()
	require.NotNil(t, viperEngine)
	viperEngine2 := getViperEngine()
	require.Equal(t, viperEngine, viperEngine2)
}

func TestViperEngineReadEnv(t *testing.T) {
	resetViperEngineInstance()
	os.Setenv(envPrefix+envSeparator+"TESTVIPERREADENV", "envValue")
	viperEngine := getViperEngine()
	value, err := viperEngine.getValue("TestViperReadEnv")
	require.Nil(t, err)
	require.Equal(t, "envValue", value)
}

func TestViperEngineReadConfigFile(t *testing.T) {
	resetViperEngineInstance()
	path := "./" + configFileName + ".yml"
	content := []byte("test: hello\n")
	ioutil.WriteFile(path, content, 0644)
	defer os.Remove(path)
	viperEngine := getViperEngine()
	value, err := viperEngine.getValue("test")
	require.Nil(t, err)
	require.Equal(t, "hello", value)
}

func TestViperEngineSetDefaultValue(t *testing.T) {
	key := "test0"
	defaultValue := "defaultTest"
	err := getViperEngine().setDefaultValue(key, defaultValue)
	require.Nil(t, err)
	value, err := getViperEngine().getValue(key)
	require.Nil(t, err)
	require.Equal(t, defaultValue, value)
}

func TestViperEngineSetValue(t *testing.T) {
	key := "test0"
	defaultValue := "defaultTest"
	settedValued := "newValue"
	getViperEngine().setDefaultValue(key, defaultValue)
	err := getViperEngine().setValue(key, settedValued)
	require.Nil(t, err)
	value, err := getViperEngine().getValue(key)
	require.Nil(t, err)
	require.Equal(t, settedValued, value)
}
func TestViperEngineSettingGetEnvKey(t *testing.T) {
	require.Equal(t, envPrefix+envSeparator+"KEY", getViperEngine().getEnvKey("key"))
}
