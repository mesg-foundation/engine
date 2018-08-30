package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// All the Log configuration keys.
const (
	LogFormat = "Log.Format"
	LogLevel  = "Log.Level"
)

func setLogDefault() {
	viper.SetDefault(LogFormat, "text")
	viper.SetDefault(LogLevel, "info")
}

func validateLog() {
	format := viper.GetString(LogFormat)
	if format != "text" && format != "json" {
		fmt.Fprintf(os.Stderr, "config: %s is not valid log format", format)
		os.Exit(1)
	}

	level := viper.GetString(LogLevel)
	if _, err := logrus.ParseLevel(level); err != nil {
		fmt.Fprintf(os.Stderr, "config: %s is not valid log level", level)
		os.Exit(1)
	}
}
