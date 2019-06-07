package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/sirupsen/logrus"
)

const (
	envPrefix = "mesg"

	serviceDBVersion   = "v2"
	executionDBVersion = "v1"
)

var (
	instance *Config
	once     sync.Once
)

// Config contains all the configuration needed.
type Config struct {
	Name string
	Path string

	Server struct {
		Address string
	}

	Log struct {
		Format      string
		ForceColors bool
		Level       string
	}

	Database struct {
		ServiceRelativePath   string
		ExecutionRelativePath string
	}

	Service ServiceConfigGroup
}

// New creates a new config with default values.
func New() (*Config, error) {
	var c Config
	c.Server.Address = ":50052"
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Log.ForceColors = false
	c.Name = "engine"
	c.Path = "/mesg"
	c.Database.ServiceRelativePath = filepath.Join("database", "services", serviceDBVersion)
	c.Database.ExecutionRelativePath = filepath.Join("database", "executions", executionDBVersion)
	c.Service = c.getServiceConfigGroup()
	return &c, nil
}

// Global returns a singleton of a Config after loaded ENV and validate the values.
func Global() (*Config, error) {
	var err error
	once.Do(func() {
		instance, err = New()
		if err != nil {
			return
		}
		if err = instance.Load(); err != nil {
			return
		}
		if err = instance.Prepare(); err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	if err := instance.Validate(); err != nil {
		return nil, err
	}
	return instance, nil
}

// Load reads config from environmental variables.
func (c *Config) Load() error {
	envconfig.MustProcess(envPrefix, c)
	return nil
}

// Prepare setups local directories or any other required thing based on config
func (c *Config) Prepare() error {
	return os.MkdirAll(c.Path, os.FileMode(0755))
}

// Validate checks values and return an error if any validation failed.
func (c *Config) Validate() error {
	if !xstrings.SliceContains([]string{"text", "json"}, c.Log.Format) {
		return fmt.Errorf("value %q is not an allowed", c.Log.Format)
	}
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return err
	}
	return nil
}
