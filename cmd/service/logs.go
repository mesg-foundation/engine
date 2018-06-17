package service

import (
	"context"
	"io"
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
	go showLogs(args[0], cmd.Flag("dependency").Value.String())
	<-utils.WaitForCancel()
}

func showLogs(serviceID string, dependency string) (err error) {
	reply, err := cli.GetService(context.Background(), &core.GetServiceRequest{
		ServiceID: serviceID,
	})
	utils.HandleError(err)
	readers, err := reply.Service.Logs(dependency)
	utils.HandleError(err)
	for _, reader := range readers {
		go func(r io.ReadCloser) {
			defer r.Close()
			stdcopy.StdCopy(os.Stdout, os.Stderr, r)
		}(reader)
	}
	return
}
