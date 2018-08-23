package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mesg-foundation/core/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	envPrefix        = "MESG"
	defaultSeparator = "."
	envSeparator     = "_"
	configFileName   = "config"
	configDir        = ".mesg"
)

// All the configuration keys.
const (
	APIServerAddress     = "Api.Server.Address"
	APIServerSocket      = "Api.Server.Socket"
	APIClientTarget      = "Api.Client.Target"
	APIServiceTargetPath = "Api.Service.TargetPath"
	APIServiceSocketPath = "Api.Service.SocketPath"
	LogFormat            = "Log.Format"
	LogLevel             = "Log.Level"
	ServicePathHost      = "Service.Path.Host"
	ServicePathDocker    = "Service.Path.Docker"
	MESGPath             = "MESG.Path"
	CoreImage            = "Core.Image"
)

var defaultDir = getDefaultDir()

func getDefaultDir() string {
	homeDir, err := homedir.Dir()
	if err == nil {
		return filepath.Join(homeDir, configDir)
	}

	u, err := user.Current()
	if err == nil {
		return filepath.Join(u.HomeDir, configDir)
	}

	return ""
}

func setDefault() {
	viper.SetDefault(MESGPath, defaultDir)

	viper.SetDefault(APIServerAddress, ":50052")
	viper.SetDefault(APIServerSocket, "/mesg/server.sock")

	viper.SetDefault(APIClientTarget, viper.GetString(APIServerAddress))

	viper.SetDefault(APIServiceSocketPath, filepath.Join(viper.GetString(MESGPath), "server.sock"))
	viper.SetDefault(APIServiceTargetPath, "/mesg/server.sock")

	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")

	viper.SetDefault(ServicePathHost, filepath.Join(viper.GetString(MESGPath), "services"))
	viper.SetDefault(ServicePathDocker, filepath.Join("/mesg", "services"))

	// Keep only the first part if Version contains space
	coreTag := strings.Split(version.Version, " ")
	viper.SetDefault(CoreImage, "mesg/core:"+coreTag[0])

}

type Config struct {
	API struct {
		Server struct {
			Address string
			Socket  string
		}
		Clinet struct {
			Target string
		}
		Service struct {
			TargetPath string
			SocketPath string
		}
	}
	Log struct {
		Format string
		Level  string
	}
	Service struct {
		Path struct {
			Host   string
			Docker string
		}
	}
	Mesg struct {
		Path string
	}
	Core struct {
		Image string
	}
}

func New() (*Config, error) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(defaultSeparator, envSeparator))
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	setDefault()
	if err := os.MkdirAll(viper.GetString(ServicePathDocker), os.ModePerm); err != nil {
		return nil, err
	}

	var c Config

	c.API.Server.Address = viper.GetString(APIServerAddress)
	c.API.Server.Socket = viper.GetString(APIServerSocket)
	c.API.Clinet.Target = viper.GetString(APIClientTarget)
	c.API.Service.TargetPath = viper.GetString(APIServiceTargetPath)
	c.API.Service.SocketPath = viper.GetString(APIServiceSocketPath)
	c.Log.Format = viper.GetString(LogFormat)
	c.Log.Level = viper.GetString(LogLevel)
	c.Service.Path.Host = viper.GetString(ServicePathHost)
	c.Service.Path.Docker = viper.GetString(ServicePathDocker)
	c.Mesg.Path = viper.GetString(MESGPath)
	c.Core.Image = viper.GetString(CoreImage)

	return &c, nil
}

func (c *Config) validate() error {
	format := viper.GetString(LogFormat)
	if format != "text" && format != "json" {
		return fmt.Errorf("config: %s is not valid log format", format)
	}

	level := viper.GetString(LogLevel)
	if _, err := logrus.ParseLevel(level); err != nil {
		return fmt.Errorf("config: %s is not valid log level", level)
	}
	return nil
}
