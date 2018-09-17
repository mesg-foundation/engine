package commands

import (
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	baseCmd

	noColor   bool
	noSpinner bool
}

func newRootCmd(e Executor) *rootCmd {
	c := &rootCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:              "mesg-core",
		Short:            "MESG Core",
		PersistentPreRun: c.persistentPreRun,
		RunE:             c.runE,
		SilenceUsage:     true,
		SilenceErrors:    true,
	})
	c.cmd.PersistentFlags().BoolVar(&c.noColor, "no-color", c.noColor, "disable colorized output")
	c.cmd.PersistentFlags().BoolVar(&c.noSpinner, "no-spinner", c.noSpinner, "disable spinners")
	c.cmd.AddCommand(
		newRootServiceCmd(e).cmd,
	)
	return c
}

func (c *rootCmd) persistentPreRun(cmd *cobra.Command, args []string) {
	if c.noColor {
		pretty.DisableColor()
	}
	if c.noSpinner {
		pretty.DisableSpinner()
	}
}

func (c *rootCmd) runE(cmd *cobra.Command, args []string) error {
	cfg, err := config.Global()
	if err != nil {
		return err
	}

	db, err := database.NewServiceDB(cfg.Database.Path)
	if err != nil {
		return err
	}

	logger.Init(cfg.Log.Format, cfg.Log.Level)

	logrus.Println("Starting Core", version.Version)

	tcpServer := &grpc.Server{
		Network:   "tcp",
		Address:   cfg.Server.Address,
		ServiceDB: db,
	}

	if err := tcpServer.Serve(); err != nil {
		return err
	}
	tcpServer.Close()
	return nil
}
