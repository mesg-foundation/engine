package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// logFormatValue represents log format flag value.
type logFormatValue string

func (v *logFormatValue) Set(value string) error {
	if value != "text" && value != "json" {
		return fmt.Errorf("%s is not valid log format", value)
	}
	*v = logFormatValue(value)
	return nil
}
func (v *logFormatValue) Type() string   { return "string" }
func (v *logFormatValue) String() string { return string(*v) }

// logLevelValue represents log level flag value.
type logLevelValue string

func (v *logLevelValue) Set(value string) error {
	if _, err := logrus.ParseLevel(value); err != nil {
		return fmt.Errorf("%s is not valid log level", value)
	}
	*v = logLevelValue(value)
	return nil
}
func (v *logLevelValue) Type() string   { return "string" }
func (v *logLevelValue) String() string { return string(*v) }
