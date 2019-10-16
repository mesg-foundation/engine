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
	"github.com/mesg-foundation/engine/x/xstrings"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
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
		*tmconfig.Config `validate:"required"`
		Path             string `validate:"required"`
	}

	Cosmos CosmosConfig
}

// CosmosConfig is the struct to hold cosmos related configs.
type CosmosConfig struct {
	Path               string          `validate:"required"`
	ChainID            string          `validate:"required"`
	GenesisTxPath      string          `validate:"required"`
	GenesisTime        time.Time       `validate:"required"`
	GenesisValidatorTx authtypes.StdTx `validate:"required" yaml:"-"`
}

// Default creates a new config with default values.
func Default() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var c Config
	c.Name = "testuser"
	c.Log.Format = "text"
	c.Log.Level = "info"
	c.Server.Address = ":50052"
	c.Path = filepath.Join(home, ".mesg", c.Name)

	c.Cosmos.ChainID = "mesg-testnet"
	c.Cosmos.GenesisTime = time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	c.Tendermint.Config = tmconfig.DefaultConfig()
	c.Tendermint.Config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	c.Tendermint.Config.P2P.AddrBookStrict = false
	c.Tendermint.Config.P2P.AllowDuplicateIP = true
	c.Tendermint.Config.Consensus.TimeoutCommit = 10 * time.Second
	c.Tendermint.Instrumentation.Prometheus = true
	c.Tendermint.Instrumentation.PrometheusListenAddr = "0.0.0.0:26660"
	return &c, nil
}

// New returns a Config after loaded yaml config and validate the values.
func New(filename string) (*Config, error) {
	c, err := Default()
	if err != nil {
		return nil, err
	}

	if err := c.Load(filename); err != nil {
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
func (c *Config) Load(filename string) error {
	if filename != "" {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(b, c); err != nil {
			return err
		}
	}

	c.Cosmos.Path = filepath.Join(c.Path, "cosmos")
	c.Cosmos.GenesisTxPath = filepath.Join(c.Path, "genesis-tx.json")
	c.Tendermint.Path = filepath.Join(c.Path, "tendermint")
	c.Tendermint.SetRoot(c.Tendermint.Path)
	c.Database.InstanceRelativePath = filepath.Join(c.Path, "database", "instance", instanceDBVersion)
	c.Database.ExecutionRelativePath = filepath.Join(c.Path, "database", "execution", executionDBVersion)
	c.Database.ProcessRelativePath = filepath.Join(c.Path, "database", "process", processDBVersion)

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

	b, err := ioutil.ReadFile(c.Cosmos.GenesisTxPath)
	if err != nil {
		return err
	}

	tx, err := decodeAuthStdTx(b)
	if err != nil {
		return err
	}

	c.Cosmos.GenesisValidatorTx = tx
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
	if err := c.Cosmos.GenesisValidatorTx.ValidateBasic(); err != nil {
		return fmt.Errorf("config.Cosmos.GenesisValidatorTx error: %w", err.Stacktrace())
	}
	return validator.New().Struct(c)
}

// decodeAuthStdTx parses string value as hex ed25519 key.
func decodeAuthStdTx(value []byte) (authtypes.StdTx, error) {
	if len(value) == 0 {
		return authtypes.StdTx{}, nil
	}

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdktypes.RegisterCodec(cdc)
	stakingtypes.RegisterCodec(cdc)

	var tx authtypes.StdTx
	if err := cdc.UnmarshalJSON(value, &tx); err != nil {
		return authtypes.StdTx{}, fmt.Errorf("unmarshal genesis validator error: %s", err)
	}
	return tx, nil
}
