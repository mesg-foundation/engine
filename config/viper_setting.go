package config

import (
	"strings"
)

type viperSetting struct {
	key string
}

func newViperSetting(key string, defaultValue string) setting {
	getViper().SetDefault(key, defaultValue)
	return &viperSetting{key: key}
}

func (s *viperSetting) setValue(value string) error {
	getViper().Set(s.key, value)
	return nil
}

func (s *viperSetting) getValue() (string, error) {
	value := getViper().GetString(s.key)
	return value, nil
}

func (s *viperSetting) getEnvKey() string {
	replacer := strings.NewReplacer(defaultSeparator, envSeparator)
	return envPrefix + envSeparator + replacer.Replace(strings.ToUpper(s.key))
}
