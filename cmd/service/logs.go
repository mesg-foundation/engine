package cmdService

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/core/cmd/utils"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/logrusorgru/aurora"
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
	serviceID := args[0]
	service, err := serviceDB.Get(serviceID)
	if err != nil {
		return
	}
	dependency := cmd.Flag("dependency").Value.String()
	readers, err := service.Logs(dependency)
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	for _, reader := range readers {
		defer reader.Close()
		go stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	}
	<-cmdUtils.WaitForCancel()
}
