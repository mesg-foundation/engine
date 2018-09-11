package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type startCmd struct {
	baseCmd

	lfv logFormatValue
	llv logLevelValue

	e RootExecutor
}

func newStartCmd(e RootExecutor) *startCmd {
	c := &startCmd{
		lfv: logFormatValue("text"),
		llv: logLevelValue("info"),
		e:   e,
	}
	c.cmd = newCommand(&cobra.Command{
		Use:     "start",
		Short:   "Start the Core",
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})

	c.cmd.Flags().Var(&c.lfv, "log-format", "log format [text|json]")
	c.cmd.Flags().Var(&c.llv, "log-level", "log level [debug|info|warn|error|fatal|panic]")
	return c
}

func (c *startCmd) preRunE(cmd *cobra.Command, args []string) error {
	// TODO: figure out how to move this to config
	viper.Set(config.LogFormat, string(c.lfv))
	viper.Set(config.LogLevel, string(c.llv))
	return nil
}

func (c *startCmd) runE(cmd *cobra.Command, args []string) error {
	var err error
	pretty.Progress("Starting Core...", func() { err = c.e.Start() })
	if err != nil {
		return err
	}

	fmt.Println(pretty.Success("Core is running"))
	return nil
}
