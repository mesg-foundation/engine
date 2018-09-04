package config

import (
	"fmt"

	"github.com/mesg-foundation/core/x/xstrings"
)

// Config is a wrapper of Setting that has validation functionality.
type Config struct {
	key         string
	engine      engine
	validations []func(value string) error
}

type option func(*Config)

func withAllowedValues(allowedValues ...string) option {
	return withValidation(func(value string) error {
		if xstrings.SliceContains(allowedValues, value) == false {
			return fmt.Errorf("Value %q is not an allowed", value)
		}
		return nil
	})
}

func withValidation(validation func(value string) error) option {
	return func(c *Config) {
		c.validations = append(c.validations, validation)
	}
}

// New initializes a Config based on a setting and optional validation function
func new(key string, defaultValue string, e engine, options ...option) *Config {
	c := &Config{
		key:    key,
		engine: e,
	}
	for _, option := range options {
		option(c)
	}
	c.engine.setDefaultValue(key, defaultValue)
	return c
}

func (config *Config) validate(value string) error {
	for _, validation := range config.validations {
		if err := validation(value); err != nil {
			return err
		}
	}
	return nil
}

// SetValue validates and set the value to the config
func (config *Config) SetValue(value string) error {
	if err := config.validate(value); err != nil {
		return err
	}
	return config.engine.setValue(config.key, value)
}

// GetValue returns the value and an error if the validation failed.
func (config *Config) GetValue() (string, error) {
	value, err := config.engine.getValue(config.key)
	if err != nil {
		return "", err
	}
	if err := config.validate(value); err != nil {
		return "", err
	}
	return value, nil
}

// GetEnvKey returns the key to use in env
func (config *Config) GetEnvKey() string {
	return config.engine.getEnvKey(config.key)
}
