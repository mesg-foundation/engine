package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/service"

	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a service file",
	Long: `Validate a service file. Check the yml format and rules.

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.tech/service/service-file.html).`,
	Example: `mesg-core service validate
mesg-core service validate ./SERVICE_FOLDER`,
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	warnings, err := service.ValidService(defaultPath(args))
	if err == nil {
		fmt.Println(aurora.Green("Service is valid"))
		return
	}

	fmt.Println(aurora.Red("Service error").Bold())
	fmt.Println(err)
	for _, warning := range warnings {
		fmt.Println(aurora.Red("The service file contains errors:").Bold())
		fmt.Println(warning)
	}
}
