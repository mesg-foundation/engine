package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Resume run the resume command for a service
var Resume = &cobra.Command{
	Use:               "resume SERVICE",
	Short:             "Resume a service",
	Long:              "Resume a service that have been paused.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace service resume ethereum",
	Run:               resumeHandler,
	DisableAutoGenTag: true,
}

func resumeHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select the account you want to use")
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO start and when ready resume (onchan) the service
	fmt.Println("service resume called", args, account)
}

func init() {
	cmdUtils.Confirmable(Resume)
	cmdUtils.Accountable(Resume)
}
