package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViperSetValue(t *testing.T) {
	key := "test0"
	defaultValue := "defaultTest"
	settedValued := "newValue"
	vs := newViperSetting(key, defaultValue)
	vs.setValue(settedValued)
	value, _ := vs.getValue()
	require.Equal(t, settedValued, value)
}

func TestViperSettingDefaultValue(t *testing.T) {
	key := "test1"
	defaultValue := "defaultTest"
	settedValued := "newValue"
	vs := newViperSetting(key, defaultValue)
	value1, _ := vs.getValue()
	require.Equal(t, defaultValue, value1)
	vs.setValue(settedValued)
	value2, _ := vs.getValue()
	require.Equal(t, settedValued, value2)
}

func TestViperSettingGetEnvKey(t *testing.T) {
	vs := newViperSetting("test2", "defaultTest")
	require.Equal(t, envPrefix+envSeparator+"TEST2", vs.getEnvKey())
}

func TestViperSettingReadEnv(t *testing.T) {
	envValue := "envValue"
	vs := newViperSetting("test3", "defaultTest")
	os.Setenv(vs.getEnvKey(), envValue)
	value1, _ := vs.getValue()
	require.Equal(t, envValue, value1)
}
