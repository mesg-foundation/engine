package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
)

var stake float64
var duration int

// Start run the start command for a service
var Start = &cobra.Command{
	Use:               "start SERVICE_FOLDER",
	Short:             "Start a service",
	Long:              "Start a service from the published available services. You have to provide a stake value and duration.",
	Example:           `mesg-cli service start SERVICE_FOLDER`,
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func startHandler(cmd *cobra.Command, args []string) {
	loadedService, err := service.ImportFromPath(defaultPath(args))
	handleError(err)
	if loadedService.IsRunning() {
		fmt.Println(aurora.Green("Service " + loadedService.Name + " is already running"))
		return
	}
	_, err = loadedService.Start()
	handleError(err)
	fmt.Println(aurora.Green("Service " + loadedService.Name + " started"))
}
