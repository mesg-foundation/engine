package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/core/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const envPrefix = "mesg"

var (
	instance *Config
	once     sync.Once

	logformats = []string{"text", "json"}
)

type plugin struct {
	Name string
	Path string
	ID   string
}

// Config contains all the configuration needed.
type Config struct {
	Server struct {
		Address string

		Plugins []plugin
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
		instance.LoadFromEnv()
	})
	if err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}
	return instance, nil
}

// Load reads config from yaml and env.
// Note that env variables have higher precedence then yaml config.
func (c *Config) Load(file string) error {
	if err := c.LoadFromYaml(file); err != nil {
		return err
	}
	c.LoadFromEnv()
	return nil
}

// LoadFromEnv reads config from environmental variables.
func (c *Config) LoadFromEnv() {
	envconfig.MustProcess(envPrefix, c)
}

// LoadFromYaml reads config from yaml file.
func (c *Config) LoadFromYaml(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("config: can't read file %s", err)
	}
	if err := yaml.UnmarshalStrict(content, c); err != nil {
		return fmt.Errorf("config: can't parse yaml file %s", err)
	}
	return nil
}

// Validate checks values and return an error if any validation failed.
func (c *Config) Validate() error {
	if !xstrings.SliceContains(logformats, c.Log.Format) {
		return fmt.Errorf("value %q is not allowed", c.Log.Format)
	}
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return err
	}

	for _, p := range c.Server.Plugins {
		if p.Name == "" {
			return errors.New("plugin name required")
		}
		if p.Path != "" && p.ID != "" {
			return errors.Errorf("plugin %s is both path and id - only one is allowed at the same time", p.Name)
		}
	}

	return nil
}
