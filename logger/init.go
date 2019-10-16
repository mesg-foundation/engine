package logger

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Init initializes default logger. It panics on invalid format or level.
func Init(format, level string, forceColors bool) {
	switch format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: forceColors,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
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
