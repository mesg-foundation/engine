package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/mesg-foundation/engine/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"gopkg.in/yaml.v2"
)

const (
	envPrefix = "mesg"

	serviceDBVersion   = "v4"
	executionDBVersion = "v3"
	instanceDBVersion  = "v2"
	processDBVersion   = "v2"
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
		ProcessRelativePath   string
	}

	Tendermint struct {
		*tmconfig.Config
		Path string
	}

	Cosmos CosmosConfig
}

// CosmosConfig is the struct to hold cosmos related configs.
type CosmosConfig struct {
	Path        string
	ChainID     string
	GenesisTime time.Time

	GenesisValidatorTx StdTx
}

// New creates a new config with default values.
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

	c.Database.ServiceRelativePath = filepath.Join("database", "services", serviceDBVersion)
	c.Database.InstanceRelativePath = filepath.Join("database", "instance", instanceDBVersion)
	c.Database.ExecutionRelativePath = filepath.Join("database", "executions", executionDBVersion)
	c.Database.ProcessRelativePath = filepath.Join("database", "processes", processDBVersion)

	c.Tendermint.Config = tmconfig.DefaultConfig()
	c.Tendermint.Config.P2P.AddrBookStrict = false
	c.Tendermint.Config.P2P.AllowDuplicateIP = true
	c.Tendermint.Config.Consensus.TimeoutCommit = 10 * time.Second

	return &c, nil
}

// New returns a Config after loaded ENV and validate the values.
func New() (*Config, error) {
	c, err := Default()
	if err != nil {
		return nil, err
	}
	if err := c.LoadEnv(); err != nil {
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

// LoadYaml loads config from env variables.
func (c *Config) LoadEnv() error {
	if err := envconfig.Process(envPrefix, c); err != nil {
		return err
	}

	return c.load()
}

// LoadYaml loads config from yaml file.
func (c *Config) LoadYaml(file string) error {
	in, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(in, c); err != nil {
		return err
	}

	return c.load()
}

// Load reads config from environmental variables.
func (c *Config) load() error {
	c.Tendermint.Path = filepath.Join(c.Path, "tendermint")
	c.Cosmos.Path = filepath.Join(c.Path, "cosmos")

	c.Tendermint.SetRoot(c.Tendermint.Path)
	return nil
}

// Prepare setups local directories or any other required thing based on config
func (c *Config) Prepare() error {
	if err := os.MkdirAll(c.Path, os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.Tendermint.Path, "config"), os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.Tendermint.Path, "data"), os.FileMode(0755)); err != nil {
		return err
	}
	if err := os.MkdirAll(c.Cosmos.Path, os.FileMode(0755)); err != nil {
		return err
	}
	return nil
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

// StdTx is type used to parse cosmos tx value provided by envconfig.
type StdTx authtypes.StdTx

// Decode parses string value as hex ed25519 key.
func (tx *StdTx) Decode(value string) error {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdktypes.RegisterCodec(cdc)
	stakingtypes.RegisterCodec(cdc)

	if err := cdc.UnmarshalJSON([]byte(value), tx); err != nil {
		return fmt.Errorf("unmarshal genesis validator error: %s", err)
	}
	return nil
}
