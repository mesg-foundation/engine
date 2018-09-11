package main

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/core/cmd"
	"github.com/mesg-foundation/core/commands"
	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/utils/clierrors"
	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	cmd.RootCmd.Version = version.Version
	cmd.RootCmd.Short = cmd.RootCmd.Short + " " + version.Version
}

func main() {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	p := provider.New(core.NewCoreClient(connection))
	if err := commands.Build(p).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, clierrors.ErrorMessage(err))
		os.Exit(1)
	}
}
