package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Init initializes default logger. It panics on invalid format or level.
func Init(format, level string) {
	switch format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		panic(fmt.Sprintf("log: %s is not a valid format", format))
	}

	l, err := logrus.ParseLevel(level)
	if err != nil {
		panic(fmt.Sprintf("log: %s is not a valid level", level))
	}
	logrus.SetLevel(l)
}
