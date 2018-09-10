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

func TestNewEntry(t *testing.T) {
	entry := newEntry("key", "defaultValue1", newTestEngine())
	require.NotNil(t, entry)
	require.NotNil(t, entry.engine)
}

func TestEntryGetSetValue(t *testing.T) {
	entry := newEntry("key", "testValue", newTestEngine())
	value, err := entry.GetValue()
	require.Nil(t, err)
	require.Equal(t, "testValue", value)
	err = entry.SetValue("newValue")
	require.Nil(t, err)
	value, err = entry.GetValue()
	require.Nil(t, err)
	require.Equal(t, "newValue", value)
}

func TestEntrySetValueError(t *testing.T) {
	entry := newEntry("", "", &testErrorEngine{})
	err := entry.SetValue("leu")
	require.NotNil(t, err)
}

func TestEntryGetValueError(t *testing.T) {
	entry := newEntry("", "", &testErrorEngine{})
	value, err := entry.GetValue()
	require.NotNil(t, err)
	require.Equal(t, "", value)
}

func TestEntryGetEnvKey(t *testing.T) {
	entry := newEntry("key", "testValue", newTestEngine())
	envKey := entry.GetEnvKey()
	require.Equal(t, "KEY", envKey)
}

func TestValidationWithAllowedValues(t *testing.T) {
	entry := newEntry("key", "three", newTestEngine(), withAllowedValues("one", "two"))

	value, err := entry.GetValue()
	require.NotNil(t, err)
	require.Equal(t, "", value)

	err = entry.SetValue("three")
	require.NotNil(t, err)

	err = entry.SetValue("two")
	require.Nil(t, err)

	value, err = entry.GetValue()
	require.Nil(t, err)
	require.Equal(t, "two", value)
}
