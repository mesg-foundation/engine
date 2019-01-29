package config

// Default endpoints to access services. These endpoints are overritten by the build
// Use the following format for the variable name: "[Service]URL" (where Service is the name of the service capitalized)
// The service name should be the name of the directory inside `systemservices` but capitalized
// example:
// 		var (
// 			FooURL string
// 			BarURL string
// 		)
var ()

func (c *Config) initServices() {
}

// Services returns all services that the configuration package is aware of
func (c *Config) Services() []ServiceConfig {
	return []ServiceConfig{}
}
