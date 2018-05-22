package cmdService

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/spf13/cobra"
)

// Stop run the stop command for a service
var Stop = &cobra.Command{
	Use:   "stop SERVICE_ID",
	Short: "Stop a service",
	Long: `Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation]().`, // TODO: add link
	Example:           `mesg-core service stop SERVICE_ID`,
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	_, err := cli.StopService(context.Background(), &core.StopServiceRequest{
		ServiceID: args[0],
	})
	handleError(err)
	fmt.Println(aurora.Green("Service stopped"))
}
