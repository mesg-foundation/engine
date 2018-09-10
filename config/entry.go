package config

import (
	"fmt"

	"github.com/mesg-foundation/core/x/xstrings"
)

// Entry is a wrapper of Setting that has validation functionality.
type Entry struct {
	key         string
	engine      engine
	validations []func(value string) error
}

// option for configuring an entry
type option func(*Entry)

// withAllowedValues accepts a list of allowed values that the entry's value should match
func withAllowedValues(allowedValues ...string) option {
	return withValidation(func(value string) error {
		if xstrings.SliceContains(allowedValues, value) == false {
			return fmt.Errorf("Value %q is not an allowed", value)
		}
		return nil
	})
}

// withValidation accepts a generic function to validate the value
func withValidation(validation func(value string) error) option {
	return func(c *Entry) {
		c.validations = append(c.validations, validation)
	}
}

// newEntry initializes an entry based on a setting and optional validation function
func newEntry(key string, defaultValue string, e engine, options ...option) *Entry {
	c := &Entry{
		key:    key,
		engine: e,
	}
	for _, option := range options {
		option(c)
	}
	c.engine.setDefaultValue(key, defaultValue)
	return c
}

// validate checks the value against the validation functions
func (entry *Entry) validate(value string) error {
	for _, validation := range entry.validations {
		if err := validation(value); err != nil {
			return err
		}
	}
	return nil
}

// SetValue validates and set the value to the entry
func (entry *Entry) SetValue(value string) error {
	if err := entry.validate(value); err != nil {
		return err
	}
	return entry.engine.setValue(entry.key, value)
}

// GetValue returns the value and an error if the validation failed.
func (entry *Entry) GetValue() (string, error) {
	value, err := entry.engine.getValue(entry.key)
	if err != nil {
		return "", err
	}
	if err := entry.validate(value); err != nil {
		return "", err
	}
	return value, nil
}

// GetEnvKey returns the key to use in env
func (entry *Entry) GetEnvKey() string {
	return entry.engine.getEnvKey(entry.key)
}
