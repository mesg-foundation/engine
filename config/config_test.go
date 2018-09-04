package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testEngine struct {
	data map[string]string
}

func newTestEngine() *testEngine {
	return &testEngine{
		data: make(map[string]string),
	}
}

func (t *testEngine) setDefaultValue(key string, defaultValue string) error {
	t.data[key] = defaultValue
	return nil
}

func (t *testEngine) setValue(key string, value string) error {
	t.data[key] = value
	return nil
}

func (t *testEngine) getValue(key string) (string, error) {
	value, ok := t.data[key]
	if ok == false {
		return "", errors.New("no value this key")
	}
	return value, nil
}

func (t *testEngine) getEnvKey(key string) string {
	return strings.ToUpper(key)
}

type testErrorEngine struct {
}

func (t *testErrorEngine) setDefaultValue(key string, defaultValue string) error {
	return errors.New("error setDefaultValue")
}

func (t *testErrorEngine) setValue(key string, value string) error {
	return errors.New("error setValue")
}

func (t *testErrorEngine) getValue(key string) (string, error) {
	return "", errors.New("error getValue")
}

func (t *testErrorEngine) getEnvKey(key string) string {
	return ""
}

func TestNew(t *testing.T) {
	config := new("key", "defaultValue1", newTestEngine())
	require.NotNil(t, config)
	require.NotNil(t, config.engine)
}

func TestConfigGetSetValue(t *testing.T) {
	config := new("key", "testValue", newTestEngine())
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
	config := new("", "", &testErrorEngine{})
	err := config.SetValue("leu")
	require.NotNil(t, err)
}

func TestConfigGetValueError(t *testing.T) {
	config := new("", "", &testErrorEngine{})
	value, err := config.GetValue()
	require.NotNil(t, err)
	require.Equal(t, "", value)
}

func TestConfigGetEnvKey(t *testing.T) {
	config := new("key", "testValue", newTestEngine())
	envKey := config.GetEnvKey()
	require.Equal(t, "KEY", envKey)
}

func TestValidationWithAllowedValues(t *testing.T) {
	config := new("key", "three", newTestEngine(), withAllowedValues("one", "two"))

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
