package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a service file",
	Long: `Validate a service file. Check the yml format and rules.

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.tech/service/service-file.html).`,
	Example: `mesg-cli service validate
mesg-cli service validate ./SERVICE_FOLDER`,
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	path := "./"
	if len(args) > 0 {
		path = args[0]
	}
	if validateServicePath(path) {
		fmt.Println(aurora.Green("Service is valid"))
	}
}
