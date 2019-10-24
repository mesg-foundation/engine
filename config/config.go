package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/engine/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"gopkg.in/go-playground/validator.v9"
)

const (
	envPrefix = "mesg"

	executionDBVersion = "v3"
	instanceDBVersion  = "v2"
	processDBVersion   = "v2"
)

// Config contains all the configuration needed.
type Config struct {
	Name string `validate:"required"`
	Path string `validate:"required"`

	Server struct {
		Address string `validate:"required"`
	}

	Log struct {
		Format      string `validate:"required"`
		ForceColors bool
		Level       string `validate:"required"`
	}

	Database struct {
		InstanceRelativePath  string `validate:"required"`
		ExecutionRelativePath string `validate:"required"`
		ProcessRelativePath   string `validate:"required"`
	}

	Tendermint struct {
		RelativePath     string `validate:"required"`
		Config       *tmconfig.Config `validate:"required"`
	}

	Cosmos struct {
		RelativePath string `validate:"required"`
	}

	DevGenesis struct {
		AccountName     string `validate:"required"`
		AccountPassword string `validate:"required"`
		ChainID         string `validate:"required"`
	}
}

// Default creates a new config with default values.
func Default() (*Config, error) {
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

	c.Database.InstanceRelativePath = filepath.Join("database", "instance", instanceDBVersion)
	c.Database.ExecutionRelativePath = filepath.Join("database", "executions", executionDBVersion)
	c.Database.ProcessRelativePath = filepath.Join("database", "processes", processDBVersion)

	c.Tendermint.RelativePath = "tendermint"
	c.Tendermint.Config = tmconfig.DefaultConfig()
	c.Tendermint.Config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	c.Tendermint.Config.P2P.AddrBookStrict = false
	c.Tendermint.Config.P2P.AllowDuplicateIP = true
	c.Tendermint.Config.Consensus.TimeoutCommit = 10 * time.Second
	c.Tendermint.Config.Instrumentation.Prometheus = true
	c.Tendermint.Config.Instrumentation.PrometheusListenAddr = "0.0.0.0:26660"

	c.Cosmos.RelativePath = "cosmos"

	c.DevGenesis.AccountName = "engine"
	c.DevGenesis.AccountPassword = "pass"
	c.DevGenesis.ChainID = "mesg-dev-chain"

	return &c, nil
}

// New returns a Config after loaded ENV and validate the values.
func New() (*Config, error) {
	c, err := Default()
	if err != nil {
		return nil, err
	}
	if err := c.Load(); err != nil {
		return nil, err
	}
	if err := c.Prepare(); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// Load reads config from environmental variables.
func (c *Config) Load() error {
	if err := envconfig.Process(envPrefix, c); err != nil {
		return err
	}
	c.Tendermint.Config.SetRoot(filepath.Join(c.Path, c.Tendermint.RelativePath))
	return nil
}

// Prepare setups local directories or any other required thing based on config
func (c *Config) Prepare() error {
	if err := os.MkdirAll(c.Path, os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(c.Tendermint.Config.GenesisFile()), os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(c.Tendermint.Config.DBDir(), os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.Path, c.Cosmos.RelativePath), os.FileMode(0755)); err != nil {
		return err
	}
	return nil
}

// Validate checks values and return an error if any validation failed.
func (c *Config) Validate() error {
	if !xstrings.SliceContains([]string{"text", "json"}, c.Log.Format) {
		return fmt.Errorf("config.Log.Format value %q is not an allowed", c.Log.Format)
	}
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return fmt.Errorf("config.Log.Level error: %w", err)
	}
	if err := c.Tendermint.Config.ValidateBasic(); err != nil {
		return fmt.Errorf("config.Tendermint error: %w", err)
	}
	return validator.New().Struct(c)
}
