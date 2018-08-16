package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Start the MESG Core.
var Start = &cobra.Command{
	Use:               "start",
	Short:             "Start the MESG Core",
	Run:               startHandler,
	DisableAutoGenTag: true,
}

// logFormatValue represents log format flag value.
type logFormatValue string

func (v *logFormatValue) Set(value string) error {
	if value != "text" && value != "json" {
		return fmt.Errorf("%s is not valid log format", value)
	}
	viper.Set(config.LogFormat, value)
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
	viper.Set(config.LogLevel, value)
	*v = logLevelValue(value)
	return nil
}
func (v *logLevelValue) Type() string   { return "string" }
func (v *logLevelValue) String() string { return string(*v) }

func init() {
	lfv := logFormatValue("text")
	Start.Flags().Var(&lfv, "log-format", "log format [text|json]")

	llv := logLevelValue("info")
	Start.Flags().Var(&llv, "log-level", "log level [debug|info|warn|error|fatal|panic]")

	RootCmd.AddCommand(Start)
}

func startHandler(cmd *cobra.Command, args []string) {
	status, err := daemon.Status()
	utils.HandleError(err)
	if status == container.RUNNING {
		fmt.Println(aurora.Green("MESG Core is running"))
		return
	}
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Starting MESG Core..."}, func() {
		_, err = daemon.Start()
	})
	utils.HandleError(err)
	fmt.Println(aurora.Green("MESG Core is running"))
}
