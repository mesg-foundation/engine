package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:               "validate SERVICE_FILE",
	Short:             "Validate a service file. Check the yml format and rules.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service validate service.yml",
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	// TODO validate the service
	fmt.Println("service valid")
}
