package cmdService

import (
	"fmt"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/service"
)

// validate a service
func validateService(path string) (valid bool) {
	valid, warnings, err := service.ValidService(path)
	if err != nil {
		fmt.Println(aurora.Red("Service error").Bold())
		fmt.Println(err)
	}
	for _, warning := range warnings {
		fmt.Println(aurora.Red("The service file contains errors:").Bold())
		fmt.Println(warning)
	}
	return
}

// Validate and import the service
func importService(path string) (instance *service.Service) {
	if !validateService(path) {
		return
	}
	instance, err := service.ImportFromFile(filepath.Join(path, "mesg.yml"))
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	return
}
