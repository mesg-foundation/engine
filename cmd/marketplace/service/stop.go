package cmdServiceMarketPlace

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Stop run the stop command for a service
var Stop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a service",
	Long: `By stoping a service, your node will not process any other actions from this service.
/!\ This action will slash your stake if you didn't respect the duration`,
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace service stop ethereum",
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	if !cmdUtils.Confirm(cmd, "Are you sure ? Your stake may be slashed !") {
		return
	}
	// TODO take stake && stop service
	// Is it really usefull to take the stake, the node will be offline anyway and we cannot trust the client
	fmt.Println("service stop called", args)
}

func init() {
	Stop.Flags().BoolP("confirm", "c", false, "Confirm")
}
