package config

import (
	"encoding/json"

	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/service"
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

// Services return the config for all services.
func (c *Config) Services() ([]ServiceConfig, error) {
	var marketplace service.Service
	var ethwallet service.Service
	if err := json.Unmarshal([]byte(marketplaceCompiled), &marketplace); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(ethwalletCompiled), &ethwallet); err != nil {
		return nil, err
	}
	return []ServiceConfig{
		{
			Key:        "Marketplace",
			Definition: &marketplace,
			Env: map[string]string{
				"MARKETPLACE_ADDRESS": EnvMarketplaceAddress,
				"TOKEN_ADDRESS":       EnvMarketplaceToken,
				"PROVIDER_ENDPOINT":   EnvMarketplaceEndpoint,
			},
		},
		// {
		// 	Key:        "Ethwallet",
		// 	Definition: &ethwallet,
		// },
	}, nil
}
