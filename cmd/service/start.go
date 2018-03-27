package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var stake float64
var duration int

// Start run the start command for a service
var Start = &cobra.Command{
	Use:   "start SERVICE_ID",
	Short: "Start a service",
	Long:  "Start a service from the published available services. You have to provide a stake value and duration.",
	Args:  cobra.MinimumNArgs(1),
	Example: `mesg-cli service start SERVICE_ID
mesg-cli service start SERVICE_ID --stake STAKE --duration DURATION  --account ACCOUNT --confirm`,
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func startHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	if stake == 0 {
		survey.AskOne(&survey.Input{
			Message: "How much do you want to stake?",
			Help:    "More details on the stake here",
			Default: "0",
		}, &stake, nil)
	}
	if duration == 0 {
		survey.AskOne(&survey.Input{
			Message: "How long will you run this service?",
			Help:    "More details on the duration here",
			Default: "0",
		}, &duration, nil)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure to run this service and stake your tokens?") {
		return
	}
	// TODO stake && start service
	fmt.Println("Service started with success", args, stake, duration, account)
}

func init() {
	cmdUtils.Confirmable(Start)
	cmdUtils.Accountable(Start)
	Start.Flags().Float64VarP(&stake, "stake", "s", 0, "The amount to stake")
	Start.Flags().IntVarP(&duration, "duration", "d", 0, "The duration you will be running this service")
}
