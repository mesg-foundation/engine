package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Stop run the stop command for a service
var Stop = &cobra.Command{
	Use:   "stop SERVICE_ID",
	Short: "Stop a service",
	Long: `Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation](https://docs.mesg.tech/service/run/).
	`,
	Args: cobra.MinimumNArgs(1),
	Example: `mesg-cli service stop SERVICE_ID
mesg-cli service stop SERVICE_ID --account ACCOUNT --confirm`,
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	if !cmdUtils.Confirm(cmd, "Are you sure? Your stake may be slashed!") {
		return
	}
	// TODO take stake && stop service
	// Is it really usefull to take the stake, the node will be offline anyway and we cannot trust the client
	fmt.Println("Service stopped with success", args, account)
}

func init() {
	cmdUtils.Confirmable(Stop)
	cmdUtils.Accountable(Stop)
}
