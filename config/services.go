package config

import (
	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/service"
)

// Default endpoints to access services. These endpoints are overritten by the build
// Use the following format for the variable name: "[service]URL" (where service is the name of the service)
// The service name should be the name of the directory inside `systemservices`
// example:
// 		var (
// 			barURL string
// 		)
var (
	ethwalletURL   string
	marketplaceURL string
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
func (c *Config) Services() []ServiceConfig {
	return []ServiceConfig{
		{
			Key: "EthWallet",
			Definition: &service.Service{
				Sid:  "ethwallet",
				Name: "Ethereum Wallet",
				Configuration: &service.Dependency{
					Key: service.MainServiceKey,
				},
				Source: ethwalletURL,
			},
		},
		{
			Key: "Marketplace",
			Definition: &service.Service{
				Sid:  "marketplace",
				Name: "Marketplace",
				Configuration: &service.Dependency{
					Key: service.MainServiceKey,
				},
				Source: marketplaceURL,
			},
			Env: map[string]string{
				"MARKETPLACE_ADDRESS": EnvMarketplaceAddress,
				"TOKEN_ADDRESS":       EnvMarketplaceToken,
				"PROVIDER_ENDPOINT":   EnvMarketplaceEndpoint,
			},
		},
	}
}
