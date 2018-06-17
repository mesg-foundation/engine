package service

import (
	"context"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
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
	reply, err := cli.Service(context.Background(), &core.ServiceRequest{
		ServiceID: args[0],
	})
	utils.HandleError(err)
	readers, err := reply.Service.Logs(cmd.Flag("dependency").Value.String())
	utils.HandleError(err)
	for _, reader := range readers {
		defer reader.Close()
		go stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	}
	<-utils.WaitForCancel()
}
