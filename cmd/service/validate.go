package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:               "validate SERVICE_PATH",
	Short:             "Validate a service. Check the mesg.yml file for format and rules and do some additional tests about the directory",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service validate /path/to/the/service/folder",
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	if validateServicePath(args[0]) {
		fmt.Println(aurora.Green("Service is valid"))
	}
}
