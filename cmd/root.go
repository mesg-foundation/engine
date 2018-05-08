package cmd

import (
	"log"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"

	"github.com/mesg-foundation/core/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "mesg-core",
	Short:             "MESG CORE",
	Run:               rootHandler,
	DisableAutoGenTag: true,
}

func rootHandler(cmd *cobra.Command, args []string) {
	log.Println("Starting MESG daemon")
	serverTCP := api.Server{
		Network: "tcp",
		Address: viper.GetString(config.APIServerAddress),
	}
	serverSocket := api.Server{
		Network: "unix",
		Address: viper.GetString(config.APIServerSocket),
	}
	go func() {
		err := serverTCP.Serve()
		defer serverTCP.Stop()
		if err != nil {
			log.Panicln(err)
		}
	}()
	go func() {
		err := serverSocket.Serve()
		defer serverSocket.Stop()
		if err != nil {
			log.Panicln(err)
		}
	}()
	<-cmdUtils.WaitForCancel()
}
