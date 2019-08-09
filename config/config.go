package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/engine/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	envPrefix = "mesg"

	serviceDBVersion   = "v3"
	executionDBVersion = "v2"
	instanceDBVersion  = "v1"
	workflowDBVersion  = "v1"
)

var (
	_instance *Config
	once      sync.Once
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
		InstanceRelativePath  string
		ExecutionRelativePath string
		WorkflowRelativePath  string
	}

	Tendermint struct {
		Path            string
		ValidatorPubKey PubKeyEd25519
		P2P             struct {
			Seeds           string
			ExternalAddress string
		}
	}

	SystemServices []*ServiceConfig
}

// New creates a new config with default values.
func New() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var c Config
	c.Server.Address = ":50052"
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Log.ForceColors = false
	c.Name = "engine"
	c.Path = filepath.Join(home, ".mesg")
	c.Database.ServiceRelativePath = filepath.Join("database", "services", serviceDBVersion)
	c.Database.InstanceRelativePath = filepath.Join("database", "instance", instanceDBVersion)
	c.Database.ExecutionRelativePath = filepath.Join("database", "executions", executionDBVersion)
	c.Database.WorkflowRelativePath = filepath.Join("database", "workflows", workflowDBVersion)
	c.Tendermint.Path = "tendermint"
	c.setupServices()
	return &c, nil
}

// Global returns a singleton of a Config after loaded ENV and validate the values.
func Global() (*Config, error) {
	var err error
	once.Do(func() {
		_instance, err = New()
		if err != nil {
			return
		}
		if err = _instance.Load(); err != nil {
			return
		}
		if err = _instance.Prepare(); err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	if err := _instance.Validate(); err != nil {
		return nil, err
	}
	return _instance, nil
}

// Load reads config from environmental variables.
func (c *Config) Load() error {
	return envconfig.Process(envPrefix, c)
}

// Prepare setups local directories or any other required thing based on config
func (c *Config) Prepare() error {
	if err := os.MkdirAll(c.Path, os.FileMode(0755)); err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(c.Path, c.Tendermint.Path), os.FileMode(0755))
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

// PubKeyEd25519 is type used to parse value provided by envconfig.
type PubKeyEd25519 ed25519.PubKeyEd25519

func (key *PubKeyEd25519) Decode(value string) error {
	if value == "" {
		return fmt.Errorf("validator public key is empty")
	}

	dec, err := hex.DecodeString(value)
	if err != nil {
		return fmt.Errorf("validator public key decode error: %s", err)
	}

	if len(dec) != ed25519.PubKeyEd25519Size {
		return fmt.Errorf("validator public key %s has invalid size", value)
	}

	copy(key[:], dec)
	return nil
}
