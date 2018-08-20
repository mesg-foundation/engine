package service

import (
	"context"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/spf13/cobra"
)

// Logs of the core
var Logs = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of a service",
	Example: `mesg-core service logs SERVICE_ID
mesg-core service logs SERVICE_ID --dependency DEPENDENCY_NAME`,
	Run:               logsHandler,
	Args:              cobra.MinimumNArgs(1),
	DisableAutoGenTag: true,
}

func init() {
	Logs.Flags().StringP("dependency", "d", "*", "Name of the dependency to only show the logs from")
}

func logsHandler(cmd *cobra.Command, args []string) {
	closeReaders := showLogs(args[0], cmd.Flag("dependency").Value.String())
	defer closeReaders()
	<-utils.WaitForCancel()
}

func showLogs(serviceID string, dependency string) func() {
	reply, err := cli().GetService(context.Background(), &core.GetServiceRequest{
		ServiceID: serviceID,
	})
	utils.HandleError(err)
	readers, err := reply.Service.Logs(dependency)
	utils.HandleError(err)
	for _, reader := range readers {
		go stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	}
	return func() {
		for _, reader := range readers {
			reader.Close()
		}
	}
}
