package service

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Start run the start command for a service
var Start = &cobra.Command{
	Use:               "start SERVICE_ID",
	Short:             "Start a service",
	Long:              "Start a service from the published available services. You have to provide a stake value and duration.",
	Example:           `mesg-core service start SERVICE_ID`,
	Args:              cobra.MinimumNArgs(1),
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func startHandler(cmd *cobra.Command, args []string) {
	var err error
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Starting service..."}, func() {
		_, err = cli.StartService(context.Background(), &core.StartServiceRequest{
			ServiceID: args[0],
		})
	})
	utils.HandleError(err)
	fmt.Println(aurora.Green("Service is running"))
}
