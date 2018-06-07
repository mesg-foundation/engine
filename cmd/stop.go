package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Stop the core
var Stop = &cobra.Command{
	Use:               "stop",
	Short:             "Stop the core",
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.AddCommand(Stop)
}

func stopHandler(cmd *cobra.Command, args []string) {
	var err error
	cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Stopping core..."}, func() {
		err = stopServices()
		if err != nil {
			return
		}
		err = daemon.Stop()
	})
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	fmt.Println(aurora.Green("Core stopped"))
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
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, hash := range hashes {
		wg.Add(1)
		go func(serviceID string) {
			defer wg.Done()
			_, errStop := cli.StopService(context.Background(), &core.StopServiceRequest{
				ServiceID: serviceID,
			})
			mutex.Lock()
			defer mutex.Unlock()
			if errStop != nil && err == nil {
				err = errStop
			}
		}(hash)
	}
	wg.Wait()
	return
}
