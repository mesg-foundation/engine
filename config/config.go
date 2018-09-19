package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

const envPrefix = "mesg"

var (
	instance *Config
	once     sync.Once
)

// Config contains all the configuration needed.
type Config struct {
	Server struct {
		Address string
	}

	Client struct {
		Address string
	}

	Log struct {
		Format string
		Level  string
	}

	Core struct {
		Image string
		Path  string
	}

	Docker struct {
		Socket string
		Core   struct {
			Path string
		}
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
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Core.Image = "mesg/core:" + strings.Split(version.Version, " ")[0]
	c.Core.Path = filepath.Join(home, ".mesg")
	c.Docker.Core.Path = "/mesg"
	c.Docker.Socket = "/var/run/docker.sock"
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
	return os.MkdirAll(c.Core.Path, os.FileMode(0755))
}

// Validate checks values and return an error if any validation failed.
func (c *Config) Validate() error {
	if xstrings.SliceContains([]string{"text", "json"}, c.Log.Format) == false {
		return fmt.Errorf("Value %q is not an allowed", c.Log.Format)
	}
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return err
	}
	return nil
}

// DaemonEnv returns the needed environmental variable for the Daemon.
func (c *Config) DaemonEnv() map[string]string {
	return map[string]string{
		"MESG_LOG_FORMAT": c.Log.Format,
		"MESG_LOG_LEVEL":  c.Log.Level,
		"MESG_CORE_PATH":  c.Docker.Core.Path,
	}
}
