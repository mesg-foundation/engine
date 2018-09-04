package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testSetting struct {
	envKey string
	value  string
}

func (s *testSetting) setValue(value string) error {
	s.value = value
	return nil
}

func (s *testSetting) getValue() (string, error) {
	return s.value, nil
}

func (s *testSetting) getEnvKey() string {
	return s.envKey
}

type testErrorSetting struct{}

func (s *testErrorSetting) setValue(value string) error {
	return errors.New("error setValue")
}

func (s *testErrorSetting) getValue() (string, error) {
	return "", errors.New("error getValue")
}

func (s *testErrorSetting) getEnvKey() string {
	return ""
}

func TestNew(t *testing.T) {
	testSetting := &testSetting{value: "testValue"}
	config := new(testSetting)
	require.NotNil(t, config)
	require.NotNil(t, config.setting)
	require.Equal(t, testSetting, config.setting)
}

func TestConfigGetSetValue(t *testing.T) {
	config := new(&testSetting{value: "testValue"})
	value, err := config.GetValue()
	require.Nil(t, err)
	require.Equal(t, "testValue", value)
	err = config.SetValue("newValue")
	require.Nil(t, err)
	value, err = config.GetValue()
	require.Nil(t, err)
	require.Equal(t, "newValue", value)
}

func TestConfigSetValueError(t *testing.T) {
	config := new(&testErrorSetting{})
	err := config.SetValue("leu")
	require.NotNil(t, err)
}

func TestConfigGetValueError(t *testing.T) {
	config := new(&testErrorSetting{})
	value, err := config.GetValue()
	require.NotNil(t, err)
	require.Equal(t, "", value)
}

func TestConfigGetEnvKey(t *testing.T) {
	config := new(&testSetting{envKey: "testEnvKey"})
	envKey := config.GetEnvKey()
	require.Equal(t, "testEnvKey", envKey)
}

func TestValidationWithAllowedValues(t *testing.T) {
	testSetting := &testSetting{value: "three"}
	config := new(testSetting, withAllowedValues("one", "two"))

	value, err := config.GetValue()
	require.NotNil(t, err)
	require.Equal(t, "", value)

	err = config.SetValue("three")
	require.NotNil(t, err)

	err = config.SetValue("two")
	require.Nil(t, err)

	value, err = config.GetValue()
	require.Nil(t, err)
	require.Equal(t, "two", value)
}

// func TestDefaultValue(t *testing.T) {
// 	tests := []struct {
// 		config       func() *Config
// 		defaultValue string
// 	}{
// 		{APIPort, "50052"},
// 		{APIAddress, ""},
// 		{LogFormat, "text"},
// 		{LogLevel, "info"},
// 	}
// 	for _, test := range tests {
// 		value, err := test.config().GetValue()
// 		require.Nil(t, err)
// 		require.Equal(t, test.defaultValue, value)
// 	}

// 	coreValue, err := CoreImage().GetValue()
// 	require.Nil(t, err)
// 	require.Contains(t, coreValue, "mesg/core:")
// }

// func TestValidation(t *testing.T) {
// 	tests := []struct {
// 		config func() *Config
// 		value  string
// 	}{
// 		// {APIPort, "50052"},
// 		// {APIAddress, ""},
// 		{LogFormat, "notValidValue"},
// 		{LogLevel, "notValidValue"},
// 	}
// 	for _, test := range tests {
// 		getViper().Set(test.config().)
// 	}
// }
