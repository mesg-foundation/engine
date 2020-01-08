package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/go-bip39"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigFileName = "config.yml"
	envPathKey            = "MESG_PATH"
	envNameKey            = "MESG_NAME"
)

// Config contains all the configuration needed.
type Config struct {
	Name string `validate:"required" yaml:"-"`
	Path string `validate:"required" yaml:"-"`

	IpfsEndpoint string `validate:"required"`

	Server struct {
		Address string `validate:"required"`
	}

	Log struct {
		Format      string `validate:"required,oneof=json text"`
		ForceColors bool
		Level       string `validate:"required"`
	}

	Tendermint struct {
		Config       *tmconfig.Config `validate:"required"`
		RelativePath string           `validate:"required"`
	}

	Cosmos struct {
		RelativePath string `validate:"required"`
	}

	DevGenesis struct {
		ChainID string `validate:"required"`
	}

	Account struct {
		Name     string `validate:"required"`
		Password string `validate:"required"`
		Mnemonic string
	}
}

// defaultConfig creates a new config with default values.
func defaultConfig() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var c Config

	c.Name = "engine"
	c.Path = filepath.Join(home, ".mesg")

	c.IpfsEndpoint = "http://ipfs.app.mesg.com:8080/ipfs/"

	c.Server.Address = ":50052"
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Log.ForceColors = false

	c.Tendermint.RelativePath = "tendermint"
	c.Tendermint.Config = tmconfig.DefaultConfig()
	c.Tendermint.Config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	c.Tendermint.Config.P2P.AddrBookStrict = false
	c.Tendermint.Config.P2P.AllowDuplicateIP = true
	c.Tendermint.Config.Consensus.TimeoutCommit = 1 * time.Second
	c.Tendermint.Config.Instrumentation.Prometheus = true
	c.Tendermint.Config.Instrumentation.PrometheusListenAddr = "0.0.0.0:26660"

	c.Cosmos.RelativePath = "cosmos"

	c.DevGenesis.ChainID = "mesg-dev-chain"

	c.Account.Name = "engine"
	c.Account.Password = "pass"

	return &c, nil
}

// New returns a Config after loaded ENV and validate the values.
func New() (*Config, error) {
	c, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	if err := c.load(); err != nil {
		return nil, err
	}
	if err := c.prepare(); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// load reads config from environmental variables.
func (c *Config) load() error {
	if envName, ok := os.LookupEnv(envNameKey); ok {
		c.Name = envName
	}
	if envPath, ok := os.LookupEnv(envPathKey); ok {
		c.Path = envPath
	}
	configFilePath := filepath.Join(c.Path, defaultConfigFileName)
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			return err
		}
		if err := yaml.UnmarshalStrict(b, c); err != nil {
			return err
		}
	}
	c.Tendermint.Config.SetRoot(filepath.Join(c.Path, c.Tendermint.RelativePath))
	return nil
}

// prepare setups local directories or any other required thing based on config
func (c *Config) prepare() error {
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

// validate checks values and return an error if any validation failed.
func (c *Config) validate() error {
	if _, err := logrus.ParseLevel(c.Log.Level); err != nil {
		return fmt.Errorf("config.Log.Level error: %w", err)
	}
	if c.Account.Mnemonic != "" && !bip39.IsMnemonicValid(c.Account.Mnemonic) {
		return fmt.Errorf("config.Account.Mnemonic error: mnemonic is not valid")
	}
	if err := c.Tendermint.Config.ValidateBasic(); err != nil {
		return fmt.Errorf("config.Tendermint error: %w", err)
	}

	return validator.New().Struct(c)
}
