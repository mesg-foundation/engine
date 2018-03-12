package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Pause run the pause command for a service
var Pause = &cobra.Command{
	Use:               "pause",
	Short:             "Pause a service",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service pause ethereum",
	Run:               pauseHandler,
	DisableAutoGenTag: true,
}

func pauseHandler(cmd *cobra.Command, args []string) {
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO pause (onchain) and then stop the service
	fmt.Println("service pause called", args)
}

func init() {
	Pause.Flags().BoolP("confirm", "c", false, "Confirm")
}
