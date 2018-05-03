package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// All the configuration keys
const (
	APIServerNetwork     = "Api.Server.Network"
	APIServerAddress     = "Api.Server.Address"
	APIClientTarget      = "Api.Client.Target"
	APIServiceTarget     = "Api.Service.Target"
	APIServiceSocketPath = "Api.Service.SocketPath"
)

func init() {
	configDir, _ := getConfigDirectory()

	viper.SetDefault(APIServerNetwork, "unix")
	viper.SetDefault(APIServerAddress, filepath.Join(configDir, "server.sock"))

	viper.SetDefault(APIClientTarget, "unix://"+viper.GetString(APIServerAddress))

	viper.SetDefault(APIServiceSocketPath, "/mesg/server.sock")
	viper.SetDefault(APIServiceTarget, "unix://"+viper.GetString(APIServiceSocketPath))
}
