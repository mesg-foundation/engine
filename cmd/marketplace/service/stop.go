package cmdServiceMarketPlace

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Stop run the stop command for a service
var Stop = &cobra.Command{
	Use:               "stop SERVICE",
	Short:             "Stop a service",
	Long:              "Stop a service. The user will get its stake back if the stake duration is reached. Otherwise, it will only get a ratio of it.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace service stop ethereum",
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account")
	if !cmdUtils.Confirm(cmd, "Are you sure ? Your stake may be slashed !") {
		return
	}
	// TODO take stake && stop service
	// Is it really usefull to take the stake, the node will be offline anyway and we cannot trust the client
	fmt.Println("service stop called", args, account)
}

func init() {
	cmdUtils.Confirmable(Stop)
	cmdUtils.Accountable(Stop)
}
