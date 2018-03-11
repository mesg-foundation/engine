package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Resume run the resume command for a service
var Resume = &cobra.Command{
	Use:     "resume",
	Short:   "Resume a service",
	Args:    cobra.MinimumNArgs(1),
	Example: "mesg-cli service resume ethereum",
	Run:     resumeHandler,
}

func resumeHandler(cmd *cobra.Command, args []string) {
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO start and when ready resume (onchan) the service
	fmt.Println("service resume called", args)
}

func init() {
	Resume.Flags().BoolP("confirm", "c", false, "Confirm")
}
