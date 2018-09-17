package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/core/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

const envPrefix = "mesg"

var (
	instance *Config
	once     sync.Once

	logformats = []string{"text", "json"}
)

// Config contains all the configuration needed.
type Config struct {
	Server struct {
		Address string

		Plugin struct {
			Resolver string
		}
	}

	Database struct {
		Path string
	}

	Client struct {
		Address string
	}

	Log struct {
		Format string
		Level  string
	}
}

// New creates a new config with default values.
func New() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var c Config
	c.Server.Address = ":50052"
	c.Client.Address = "localhost:50052"
	c.Database.Path = filepath.Join(home, ".mesg", "db")
	c.Log.Format = "text"
	c.Log.Level = "info"
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
		instance.Load()
	})
	if err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}
	return instance, nil
}

// Load reads config from env.
// Note that env variables have higher precedence then yaml config.
func (c *Config) Load() {
	envconfig.MustProcess(envPrefix, c)
}

// Validate checks values and return an error if any validation failed.
func (c *Config) Validate() error {
	if !xstrings.SliceContains(logformats, c.Log.Format) {
		return fmt.Errorf("value %q is not allowed", c.Log.Format)
	}
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return err
	}
	return nil
}
