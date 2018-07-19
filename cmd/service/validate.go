package service

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/spf13/cobra"
)

// Validate a service
var Validate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a service file",
	Long: `Validate a service file. Check the yml format and rules.

All the definitions of the service file can be found in the page [Service File from the documentation](https://docs.mesg.com/service/service-file).`,
	Example: `mesg-core service validate
mesg-core service validate ./SERVICE_FOLDER`,
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	validation, err := importer.Validate(defaultPath(args))
	utils.HandleError(err)

	validateServiceFile(validation)
	validateDockerfile(validation)
	validateSummary(validation)
}

func validateServiceFile(validation *importer.ValidationResult) {
	if validation.ServiceFileExist == false {
		fmt.Printf("%s File 'mesg.yml' does not exist\n", aurora.Red("⨯"))
		return
	}

	if len(validation.ServiceFileWarnings) > 0 {
		fmt.Printf("%s File 'mesg.yml' is not valid. See documentation: %s\n", aurora.Red("⨯"), "https://docs.mesg.com/service/service-file")
		for _, warning := range validation.ServiceFileWarnings {
			fmt.Printf("  - %s\n", warning)
		}
	} else {
		fmt.Printf("%s File 'mesg.yml' is valid\n", aurora.Green("✔"))
	}
}

func validateDockerfile(validation *importer.ValidationResult) {
	if validation.DockerfileExist {
		fmt.Printf("%s Dockerfile exists\n", aurora.Green("✔"))
	} else {
		fmt.Printf("%s Dockerfile does not exist\n", aurora.Red("⨯"))
	}
}

func validateSummary(validation *importer.ValidationResult) {
	if validation.IsValid() {
		fmt.Println(aurora.Green("Service is valid"))
	} else {
		fmt.Println(aurora.Red("Service is not valid"))
	}
}
