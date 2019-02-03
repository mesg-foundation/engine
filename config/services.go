package config

// Default endpoints to access services. These endpoints are overritten by the build
// Use the following format for the variable name: "[Service]URL" (where Service is the name of the service capitalized)
// The service name should be the name of the directory inside `systemservices` but capitalized
// example:
// 		var (
// 			BarURL string
// 		)
var (
	EthwalletURL string
)

// ServiceConfig contains information related to services that the config knows about
type ServiceConfig struct {
	URL  string
	Env  map[string]string
	Sid  string
	Hash string
}

// ServiceConfigWithKey contains information related to services that the config knows about and their key
type ServiceConfigWithKey struct {
	*ServiceConfig
	Key string
}

// ServiceConfigGroup is the struct that contains all the `ServiceConfig` related to the config
type ServiceConfigGroup struct {
	Ethwallet ServiceConfig
}

// getServiceConfigGroup return the config for all services.
// DO NOT USE STRING HERE but variable defined on top of this file `XxxURL` for config injection
func (c *Config) getServiceConfigGroup() ServiceConfigGroup {
	return ServiceConfigGroup{
		Ethwallet: ServiceConfig{URL: EthwalletURL},
	}
}

// Services returns all services that the configuration package is aware of
func (c *Config) Services() []*ServiceConfigWithKey {
	return []*ServiceConfigWithKey{
		{&c.Service.Ethwallet, "Ethwallet"},
	}
}
