package config

import (
	"strings"
)

type viperSetting struct {
	key string
}

func newViperSetting(key string, defaultValue string) Setting {
	getViper().SetDefault(key, defaultValue)
	return &viperSetting{key: key}
}

func (s *viperSetting) GetValue() string {
	return getViper().GetString(s.key)
}

func (s *viperSetting) GetEnvKey() string {
	replacer := strings.NewReplacer(defaultSeparator, envSeparator)
	return envPrefix + envSeparator + replacer.Replace(strings.ToUpper(s.key))
}
