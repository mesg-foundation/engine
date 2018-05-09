package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Stop run the stop command for a service
var Stop = &cobra.Command{
	Use:   "stop SERVICE_FOLDER",
	Short: "Stop a service",
	Long: `Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation]().`, // TODO: add link
	Example:           `mesg-cli service stop SERVICE_FOLDER`,
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	service := loadService(defaultPath(args))
	stopService(service)
	fmt.Println(aurora.Green("Service " + service.Name + " stopped"))
}
