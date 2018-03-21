package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Pause run the pause command for a service
var Pause = &cobra.Command{
	Use:               "pause SERVICE",
	Short:             "Pause a service",
	Long:              "Pause a service. The user will not get its stake back but it will also not lost it. Should always pause services before quitting the CLI, otherwise the user may loss its stake.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace service pause ethereum",
	Run:               pauseHandler,
	DisableAutoGenTag: true,
}

func pauseHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select the account you want to use")
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO pause (onchain) and then stop the service
	fmt.Println("service pause called", args, account)
}

func init() {
	cmdUtils.Confirmable(Pause)
	cmdUtils.Accountable(Pause)
}
