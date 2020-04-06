package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"gopkg.in/yaml.v2"
)

const (
	// CosmosBech32MainPrefix defines the main Bech32 prefix.
	CosmosBech32MainPrefix = "mesgtest"

	// CosmosCoinType is the mesg registered coin type from https://github.com/satoshilabs/slips/blob/master/slip-0044.md.
	CosmosCoinType = uint32(470)
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

	DefaultExecutionPrice string `validate:"required"`

	// AuthorizedPubKeys are the bech32 public key of the accounts that are authorized to call the gRPC Admin API.
	AuthorizedPubKeys []string `validate:"dive,required,bech32accpubkey" yaml:"authorized_pubkeys"`

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

		// Minimum gas prices for transactions.
		MinGasPrices string `validate:"required,deccoins"`

		// Token name to use in the staking module.
		StakeTokenDenom string `validate:"required"`

		// Power reduction between the staking token and the voting power on tendermint.
		PowerReduction int64 `validate:"required"`
	}

	DevGenesis struct {
		ChainID                 string `validate:"required"`
		InitialBalances         string `validate:"required,coins"`
		ValidatorDelegationCoin string `validate:"required,coin"`
	}

	Account struct {
		Name     string `validate:"required"`
		Password string `validate:"required"`
		Number   uint32
		Index    uint32
		Mnemonic string `validate:"omitempty,mnemonic"`
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

	c.DefaultExecutionPrice = "10000atto" // /x/execution/internal/types/params.go#DefaultMinPrice

	c.Server.Address = ":50052"
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Log.ForceColors = false

	c.Tendermint.RelativePath = "tendermint"
	c.Tendermint.Config = tmconfig.DefaultConfig()
	c.Tendermint.Config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	c.Tendermint.Config.RPC.MaxSubscriptionsPerClient = 100
	c.Tendermint.Config.P2P.AddrBookStrict = false
	c.Tendermint.Config.P2P.AllowDuplicateIP = true
	c.Tendermint.Config.Consensus.TimeoutCommit = 5 * time.Second
	c.Tendermint.Config.Instrumentation.Prometheus = true
	c.Tendermint.Config.Instrumentation.PrometheusListenAddr = "0.0.0.0:26660"
	c.Tendermint.Config.TxIndex.IndexAllKeys = true

	c.Cosmos.RelativePath = "cosmos"
	c.Cosmos.MinGasPrices = "1.0atto"
	c.Cosmos.StakeTokenDenom = "atto"
	c.Cosmos.PowerReduction = 18

	c.DevGenesis.ChainID = "mesg-dev-chain"
	c.DevGenesis.InitialBalances = "250000000000000000000000000atto"       // 250 000 000 * 10^18
	c.DevGenesis.ValidatorDelegationCoin = "1000000000000000000000000atto" // 1 000 000 * 10^18

	c.Account.Name = "engine"
	c.Account.Password = "pass"
	c.Account.Number = uint32(0)
	c.Account.Index = uint32(0)

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
		return fmt.Errorf("config log.level error: %w", err)
	}
	if err := c.Tendermint.Config.ValidateBasic(); err != nil {
		return fmt.Errorf("config tendermint error: %w", err)
	}
	return xvalidator.Struct(c)
}
