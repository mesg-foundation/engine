package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViperSettingDefaultValue(t *testing.T) {
	key := "test1"
	defaultValue := "defaultTest"
	settedValued := "newValue"
	vs := newViperSetting(key, defaultValue)
	require.Equal(t, defaultValue, vs.GetValue())
	getViper().Set(key, settedValued)
	require.Equal(t, settedValued, vs.GetValue())
}

func TestViperSettingGetEnvKey(t *testing.T) {
	vs := newViperSetting("test2", "defaultTest")
	require.Equal(t, envPrefix+envSeparator+"TEST2", vs.GetEnvKey())
}

func TestViperSettingReadEnv(t *testing.T) {
	envValue := "envValue"
	vs := newViperSetting("test3", "defaultTest")
	os.Setenv(vs.GetEnvKey(), envValue)
	require.Equal(t, envValue, vs.GetValue())
}
