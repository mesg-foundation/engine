package daemon

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Stop the daemon
var Stop = &cobra.Command{
	Use:               "stop",
	Short:             "Stop the daemon",
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	err := stopServices()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	err = daemon.Stop()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	fmt.Println(aurora.Green("Daemon is stopping"))
}

func getCli() (cli core.CoreClient, err error) {
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	if err != nil {
		return
	}
	cli = core.NewCoreClient(connection)
	return
}

func stopServices() (err error) {
	cli, err := getCli()
	if err != nil {
		return
	}
	hashes, err := service.List()
	if err != nil {
		return
	}
	for _, hash := range hashes {
		_, err := cli.StopService(context.Background(), &core.StopServiceRequest{
			ServiceID: hash,
		})
		if err != nil {
			break
		}
	}
	return
}
