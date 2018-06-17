package service

import (
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/cmd/utils"
	serviceDB "github.com/mesg-foundation/core/database/services"
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
	service, err := serviceDB.Get(serviceID)
	utils.HandleError(err)
	readers, err := service.Logs(dependency)
	utils.HandleError(err)
	for _, reader := range readers {
		defer reader.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	}
	return
}
