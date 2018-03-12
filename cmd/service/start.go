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
	Use:               "start",
	Short:             "Start a service",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service start --stake 100 --duration 10 ethereum",
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func startHandler(cmd *cobra.Command, args []string) {
	if stake == 0 {
		survey.AskOne(&survey.Input{
			Message: "How much do you want to stake (MESG) ?",
			Help:    "More details on the stake here",
			Default: "0",
		}, &stake, nil)
	}
	if duration == 0 {
		survey.AskOne(&survey.Input{
			Message: "How long will you run this service (hours) ?",
			Help:    "More details on the duration here",
			Default: "0",
		}, &duration, nil)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure to run this service and stake your tokens ?") {
		return
	}
	// TODO stake && start service
	fmt.Println("service start called", args, stake, duration)
}

func init() {
	Start.Flags().BoolP("confirm", "c", false, "Confirm")
	Start.Flags().Float64VarP(&stake, "stake", "s", 0, "The number of MESG to put on stake")
	Start.Flags().IntVarP(&duration, "duration", "d", 0, "The amount of time you will be running this/those service(s) for (in hours)")
}
