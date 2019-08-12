package config

import (
	"encoding/json"

	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/service"
)

// Default compiled version of the service. These compiled versions are overritten by the build
// Use the following format for the variable name: "[service]Compiled" (where service is the name of the service)
// The service name should be the name of the directory inside `systemservices`
// example:
// 		var (
// 			barCompiled string
// 		)
var (
	ethwalletCompiled   string
	marketplaceCompiled string
)

// Env to override on the system services
var (
	EnvMarketplaceAddress  string
	EnvMarketplaceToken    string
	EnvMarketplaceEndpoint string
)

// ServiceConfig contains information related to services that the config knows about
type ServiceConfig struct {
	Key        string
	Env        map[string]string
	Definition *service.Service
	Instance   *instance.Instance
}

// setupServices initialize all services for this config
func (c *Config) setupServices() error {
	if marketplaceCompiled != "" {
		marketplace, err := c.createServiceConfig("Marketplace", marketplaceCompiled, map[string]string{
			"MARKETPLACE_ADDRESS": EnvMarketplaceAddress,
			"TOKEN_ADDRESS":       EnvMarketplaceToken,
			"PROVIDER_ENDPOINT":   EnvMarketplaceEndpoint,
		})
		if err != nil {
			return err
		}
		c.SystemServices = append(c.SystemServices, marketplace)
	}
	if ethwalletCompiled != "" {
		ethwallet, err := c.createServiceConfig("EthWallet", ethwalletCompiled, nil)
		if err != nil {
			return err
		}
		c.SystemServices = append(c.SystemServices, ethwallet)
	}
	return nil
}

func (c *Config) createServiceConfig(key string, compilatedJSON string, env map[string]string) (*ServiceConfig, error) {
	var srv service.Service
	if err := json.Unmarshal([]byte(compilatedJSON), &srv); err != nil {
		return nil, err
	}
	srv.Configuration.Key = service.MainServiceKey
	return &ServiceConfig{
		Key:        key,
		Definition: &srv,
		Env:        env,
	}, nil
}
