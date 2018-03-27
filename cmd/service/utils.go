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
		fmt.Println(aurora.Red(err))
	}
	for _, warning := range warnings {
		fmt.Println(aurora.Brown(warning))
	}
	return
}

func importService(path string) (instance *service.Service, err error) {
	instance, err = service.ImportFromFile(filepath.Join(path, "mesg.yml"))
	return
}
