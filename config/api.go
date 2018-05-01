package config

import "github.com/spf13/viper"

var Api *api

type api struct {
	Server ServerType
	Client ClientType
}

type ServerType struct{}

func (s *ServerType) Network() string {
	return viper.GetString("Api.Server.Network")
}

func (s *ServerType) Address() string {
	return viper.GetString("Api.Server.Address")
}

type ClientType struct{}

func (c *ClientType) Target() string {
	return viper.GetString("Api.Client.Target")
}

func init() {
	viper.SetDefault("Api.Server.Network", "unix")
	viper.SetDefault("Api.Server.Address", "server.sock")

	viper.SetDefault("Api.Client.Target", "localhost:50052")

	Api = new(api)
}
